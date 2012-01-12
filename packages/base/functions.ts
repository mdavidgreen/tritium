
@func export(Text %key, Text %value) {
  export(%key) {
    set(%value)
    yield()
  }
}

@func regexp(Text %exp) {
  regexp(%exp, "")
}

@func asset(Text %name) {
  concat($asset_host, %name) {
    yield()
  }
}

@func bm(Text %name) {
  log(concat(%name, ": ", 
    time() {
      yield()
    }, "s"))
}


@func else() {
  with(/.?/) {
    yield()
  }
}

@func match(Text %target, Text %comparitor) {
  match(%target) {
    with(%comparitor) {
      yield()
    }
  }
}

@func match(Text %target, Regexp %comparitor) {
  match(%target) {
    with(%comparitor) {
      yield()
    }
  }
}

@func Text.clear() {
  set("") {
    yield()
  }
}

@func Text.text() {
  this() {
    yield()
  }
}


@func Text.replace(Regexp %search, Text %with) {
  replace(%search) {
    set(%with)
    yield()
  }
}

@func Text.replace(Text %search, Text %with) {
  replace(%search) {
    set(%with)
  }
}


@func Text.rewrite(Text %type) {
  # Wrote this when the parser was a little funky. This can be much
  # prettier with some nice lvar usage
  replace(regexp(var(concat(concat("rewrite_", %type), "_matcher")))) {
    set(var(concat(concat("rewrite_", %type), "_replacement")))
  }
}