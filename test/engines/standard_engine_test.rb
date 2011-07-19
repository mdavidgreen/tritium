#ENV["SCRIPT"] = "prepend"

require 'minitest/autorun'
require_relative '../../lib/tritium/engines/standard/engine'
require_relative 'engine_tests'

class StandardEngineTest < MiniTest::Unit::TestCase
  include Tritium::Engines
  #include EngineTests
  
  def engine_class
    Standard
  end

  def test_no_script
    engine = Standard.new("")
    input = "hi"
    result = engine.run("hi")
    assert_equal input, result
  end
  
  def test_simple_set
    engine = Standard.new("set('hi')")
    result = engine.run("world")
    assert_equal "hi", result
  end
  
  def test_html_parsing
    engine = Standard.new("html()")
    result = engine.run("")
    expected = "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\" \"http://www.w3.org/TR/REC-html40/loose.dtd\">\n\n"
    assert_equal expected, result
  end
end