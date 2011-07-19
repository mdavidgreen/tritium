require 'minitest/autorun'
require_relative '../../lib/tritium/parser/parser'
ENV["TEST"] = "true"

class ParserTest < MiniTest::Unit::TestCase
  include Tritium::Parser
  
  # ============== HELPER METHODS =================
  def path
    File.join(File.dirname(__FILE__), "scripts")
  end
  
  def script_path(script_name)
    File.join(path, script_name)
  end
  
  def read_script(script_name)
    File.read(script_path(script_name))
  end
  
  def build_parser(script_name)
    script_string = read_script(script_name)
    parser = Parser.new(script_string, filename: script_name, path: path)
    [script_string, parser, parser.parse]
  end
  
  # ============== TESTS!!!!!!!!! ====================
  
  def self.test_version_folder(version_dir)
    version = File.basename(version_dir)
     puts "Testing parser against version #{version}"
    Dir[version_dir + "/scripts/*.ts"].each do |script_file_name|
      test_name = File.basename(script_file_name, ".ts")
      if ENV["SCRIPT"].nil? || test_name == ENV["SCRIPT"]
        # Writes a method that simply calls run_file with its name 
        eval "def test_parser_#{version}_#{test_name}_script; run_test '#{version_dir}', '#{test_name}'; end"
      end
    end
  end

  [2].each do |ver|
    test_version_folder(File.join(File.dirname(__FILE__), "../functional/v#{ver}"))
  end
  
  def run_test(version_dir, test_name)
    parser = Parser.new(filename: test_name + ".ts", path: File.join(version_dir, "scripts"))
    parser.parse
    assert true # If we didn't error, its a positive assertion
  end

  def test_parser
    script_string, parser, output = build_parser("false-negatives.ts")
    assert_equal read_script("reference-output.ts"), output.to_s
  end
  
  def test_var_expansion
    script_string, parser, output = build_parser("var.ts")
    assert_equal read_script("var_output.ts"), output.to_s
  end
  
  def test_multiline_comment_unclosed
    begin
      script_string, parser, output = build_parser("unclosed_comment.ts")
      assert false, "Should throw error on unclosed comment"
    rescue Tritium::Parser::Parser::SyntaxError
      assert true
    end
  end
  
  def test_stop_second_parse_call
    script_string, parser, output = build_parser("basic.ts")
    
    begin
      parser.parse
      assert false, "Can't call Parser#parse twice. Should have thrown error"
    rescue Exception
      assert true
    end
  end
  
  def test_add_class_expansion
    script_string, parser, output = build_parser("add_class.ts")
    assert_equal read_script("add_class_output.ts"), output.to_s
  end
  
  def test_insert_positionals
    script_string, parser, output = build_parser("insert.ts")
    assert true # no errors
  end
  
  def test_error_handling
    build_parser("invalid.ts")
    assert false, "Should have failed"
  rescue Exception => e
    assert e.any?
    assert e.size > 2
  end
  
  def test_scopes
    script_string, parser, root = build_parser("add_class.ts")
    assert_equal "Text", root.scope.name
    expansion = root.statements.first
    var = expansion.statements.first
    assert_equal root.scope, var.scope
    assert_equal "Text", var.opens.name
    html = root.statements.last
    assert_equal root.scope, html.scope
    assert_equal "XMLNode", html.opens.name
    child = html.statements.first
    assert_equal html, child.parent
    assert !html.base?, "html() is NOT a Base method"
    assert var.base?, "var() should be a base method #{var.class} #{var.spec.scopes.inspect}"
    assert_equal "XMLNode", child.scope.name
  end
end
