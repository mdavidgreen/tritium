package shark

import(
	"strings"
	"os"
	"libxml"
	"fmt"
	xml "libxml/tree"
	tp "tritium/proto"
	"libxml/xpath"
	"rubex"
)

func (ctx *Ctx) runBuiltIn(fun *Function, scope *Scope, ins *tp.Instruction, yieldBlock *tp.Instruction, args []interface{}) (returnValue interface{}) {
	switch fun.Name {
	case "this":
		returnValue = scope.Value
	case "yield": 
		if yieldBlock != nil {
			returnValue = ctx.runChildren(scope, yieldBlock, nil)
			yieldBlock = nil
		}

	case "var.Text":
		val := ctx.Env[args[0].(string)]
		ts := &Scope{Value: val}
		ctx.runChildren(ts, ins, yieldBlock)
		returnValue = ts.Value
		ctx.Env[args[0].(string)] = returnValue.(string)
	case "var.Text.Text":
		ctx.Env[args[0].(string)] = args[1].(string)
		returnValue = args[1].(string)
	case "match.Text":
		// Setup stacks
		ctx.MatchStack = append(ctx.MatchStack, args[0].(string))
		ctx.MatchShouldContinue = append(ctx.MatchShouldContinue, true)
	
		// Run children
		ctx.runChildren(scope, ins, yieldBlock)
	
		if ctx.matchShouldContinue() {
			returnValue = "false"
		} else {
			returnValue = "true"
		}
	
		// Clear
		ctx.MatchShouldContinue = ctx.MatchShouldContinue[:len(ctx.MatchShouldContinue)-1]
		ctx.MatchStack = ctx.MatchStack[:len(ctx.MatchStack)-1]
	case "with.Text":
		returnValue = "false"
		if ctx.matchShouldContinue() {
			if args[0].(string) == ctx.matchTarget() {
				ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
				ctx.runChildren(scope, ins, yieldBlock)
				returnValue = "true"
			}
		}
	case "with.Regexp":
		returnValue = "false"
		if ctx.matchShouldContinue() {
			//println(matcher.MatchAgainst, matchWith)
			if (args[0].(*rubex.Regexp)).Match([]uint8(ctx.matchTarget())) {
				ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
				ctx.runChildren(scope, ins, yieldBlock)
				returnValue = "true"
			}
		}
	case "not.Text":
		returnValue = "false"
		if ctx.matchShouldContinue() {
			if args[0].(string) != ctx.matchTarget() {
				ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
				ctx.runChildren(scope, ins, yieldBlock)
				returnValue = "true"
			}
		}
	case "not.Regexp":
		returnValue = "false"
		if ctx.matchShouldContinue() {
			//println(matcher.MatchAgainst, matchWith)
			if !(args[0].(*rubex.Regexp)).Match([]uint8(ctx.matchTarget())) {
				ctx.MatchShouldContinue[len(ctx.MatchShouldContinue)-1] = false
				ctx.runChildren(scope, ins, yieldBlock)
				returnValue = "true"
			}
		}
	case "regexp.Text.Text":
		mode := rubex.ONIG_OPTION_DEFAULT
		if strings.Index(args[1].(string), "i") >= 0 {
			mode = rubex.ONIG_OPTION_IGNORECASE
		}
		if strings.Index(args[1].(string), "m") >= 0 {
			mode = rubex.ONIG_OPTION_MULTILINE
		}
		var err os.Error
		returnValue, err = rubex.NewRegexp(args[0].(string), mode)
		if err != nil {
			panic("Invalid regexp")
		}
	case "export.Text":
		val := make([]string, 2)
		val[0] = args[0].(string)
		ts := &Scope{Value:""}
		ctx.runChildren(ts, ins, yieldBlock)
		val[1] = ts.Value.(string)
		ctx.Exports = append(ctx.Exports, val)
	case "log.Text":
		ctx.Logs = append(ctx.Logs, args[0].(string))

	// ATOMIC FUNCTIONS
	case "concat.Text.Text":
		//println("Concat:", args[0].(string), "+", args[1].(string))
		returnValue = args[0].(string) + args[1].(string)
	case "concat.Text.Text.Text": //REMOVE
		returnValue = args[0].(string) + args[1].(string) + args[2].(string)
	case "downcase.Text":
		returnValue = strings.ToLower(args[0].(string))
		return
	case "upcase.Text":
		returnValue = strings.ToUpper(args[0].(string))
		return
	case "index.XMLNode":
		returnValue = fmt.Sprintf("%d", scope.Index + 1)
	
	// TEXT FUNCTIONS
	case "set.Text":
		scope.Value = args[0]
	case "append.Text":
		scope.Value = scope.Value.(string) + args[0].(string)
	case "prepend.Text":
		scope.Value = args[0].(string) + scope.Value.(string)
	case "replace.Text":
		ts := &Scope{Value:""}
		ctx.runChildren(ts, ins, yieldBlock)
		scope.Value = strings.Replace(scope.Value.(string), args[0].(string), ts.Value.(string), -1)
	case "replace.Regexp":
		regexp := args[0].(*rubex.Regexp)
		scope.Value = regexp.GsubFunc(scope.Value.(string), func(match string, captures map[string]string) string {
			usesGlobal := (ctx.Env["use_global_replace_vars"] == "true")
		
			for name, capture := range captures {
				if usesGlobal {
					//println("setting $", name, "to", capture)
					ctx.Env[name] = capture
				}
				ctx.LocalVar[name] = capture
			}

			replacementScope := &Scope{Value:match}
			ctx.runChildren(replacementScope, ins, yieldBlock)
			//println(ins.String())
		
			//println("Replacement:", replacementScope.Value.(string))
			innerReplacer := rubex.MustCompile(`[\\\$]([\d])`)
			return innerReplacer.GsubFunc(replacementScope.Value.(string), func(_ string, numeric_captures map[string]string) string {
				capture := numeric_captures["1"]
				var val string
				if usesGlobal {
					val = ctx.Env[capture]
				} else {
					val = ctx.LocalVar[capture].(string)
				}
				return val
		    })
		})
		returnValue = scope.Value

	// XML FUNCTIONS
	case "xml":
		doc := libxml.XmlParseString(scope.Value.(string))
		defer doc.Free()
		ns := &Scope{Value:doc}
		ctx.runChildren(ns, ins, yieldBlock)
		scope.Value = doc.String()
		returnValue = scope.Value
	case "html":
		doc := libxml.HtmlParseString(scope.Value.(string))
		defer doc.Free()
		ns := &Scope{Value:doc}
		ctx.runChildren(ns, ins, yieldBlock)
		scope.Value = doc.DumpHTML()
		returnValue = scope.Value
	case "html_fragment":
		doc := libxml.HtmlParseFragment(scope.Value.(string))
		defer doc.Free()
		ns := &Scope{Value: doc.RootElement()}
		ctx.runChildren(ns, ins, yieldBlock)
		scope.Value = ns.Value.(xml.Node).Content()
		returnValue = scope.Value
	case "select.Text":
		// TODO reuse XPath object
		node := scope.Value.(xml.Node)
		xpCtx := xpath.NewXPath(node.Doc())
		xpath := xpath.CompileXPath(args[0].(string))
		//xpath := ctx.XPath(args[0].(string))
		nodeSet := xpCtx.SearchByCompiledXPath(node, xpath).Slice()
		defer xpCtx.Free()
		if len(nodeSet) == 0 {
			returnValue = "false"
		} else {
			returnValue = "true"
		}

		for index, node := range(nodeSet) {
			if (node != nil) && node.IsLinked() && node.IsValid() {
				ns := &Scope{Value: node, Index: index}
				ctx.runChildren(ns, ins, yieldBlock)
			}
		}
	case "position.Text":
		returnValue = Positions[args[0].(string)]
	
	// SHARED NODE FUNCTIONS
	case "remove":
		scope.Value.(xml.Node).Remove()
	case "inner", "value", "inner_text", "text":
		node := scope.Value.(xml.Node)
		ts := &Scope{Value:node.Content()}
		ctx.runChildren(ts, ins, yieldBlock)
		val := ts.Value.(string)
		if node.IsLinked() {
			node.SetContent(val)
		}
		returnValue = val
	case "name":
		node := scope.Value.(xml.Node)
		ts := &Scope{Value:node.Name()}
		ctx.runChildren(ts, ins, yieldBlock)
		node.SetName(ts.Value.(string))
		returnValue = ts.Value.(string)
	case "dup":
		node := scope.Value.(xml.Node)
		newNode := node.Duplicate()
		MoveFunc(newNode, node, AFTER)
		ns := &Scope{Value:newNode}
		ctx.runChildren(ns, ins, yieldBlock)
	case "fetch.Text":
		searchNode := scope.Value.(xml.Node)
		xPathObj := xpath.NewXPath(searchNode.Doc())
		nodeSet := xPathObj.Search(searchNode, args[0].(string))
		if nodeSet.Size() > 0 {
			node := nodeSet.First()
			attr, ok := node.(*xml.Attribute)
			if ok {
				returnValue = attr.Content()
			} else {
				returnValue = node.String()
			}
		}
		xPathObj.Free()

	// LIBXML FUNCTIONS
	case "insert_at.Position.Text":
		node := scope.Value.(xml.Node)
		position := args[0].(Position)
		tagName := args[1].(string)
		element := node.Doc().NewElement(tagName)
		MoveFunc(element, node, position)
		ns := &Scope{Value: element}
		ctx.runChildren(ns, ins, yieldBlock)
	case "inject_at.Position.Text":
		node := scope.Value.(xml.Node)
		position := args[0].(Position)
		nodeSet := node.Doc().ParseHtmlFragment(args[1].(string))
		for _, newNode := range(nodeSet) {
			MoveFunc(newNode, node, position)
		}
		if len(nodeSet) > 0 {
			element, ok := nodeSet[0].(*xml.Element)
			if ok {
				// successfully ran scope
				returnValue = "true"
				ns := &Scope{Value: element}
				ctx.runChildren(ns, ins, yieldBlock)
			}
		} else {
			returnValue = "false"
		}
	case "move.XMLNode.XMLNode.Position", "move.Node.Node.Position":
		//for name, value := range(ctx.LocalVar) {
		//	println(name, ":", value)
		//}
		MoveFunc(args[0].(xml.Node), args[1].(xml.Node), args[2].(Position))
	case "wrap_text_children.Text":
		returnValue = "false"
		child := scope.Value.(xml.Node).First()
		index := 0
		tagName := args[0].(string)
		for child != nil {
			text, ok := child.(*xml.Text)
			childNext := child.Next()
			if ok {
				returnValue = "true"
				wrap := text.Wrap(tagName)
				ns := &Scope{wrap, index}
				ctx.runChildren(ns, ins, yieldBlock)
				index++
			}
			child = childNext
		}

	// ATTRIBUTE FUNCTIONS
	case "attribute.Text":
		node := scope.Value.(xml.Node)
		attr, _ := node.Attribute(args[0].(string))
		as := &Scope{Value:attr}
		ctx.runChildren(as, ins, yieldBlock)
		if attr.IsLinked() && (attr.Content() == "") {
			attr.Remove()
		}
		if !attr.IsLinked() {
			attr.Free()
		}
		returnValue = "true"
	
	default:
		ctx.Log.Error("Must implement " + fun.Name)
	}
	return
}