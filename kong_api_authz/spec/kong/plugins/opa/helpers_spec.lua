describe("opa:helpers", function()
  local helpers

  setup(function()
    -- add mocked modules to path
    package.path = package.path..";spec/?.lua"
    -- override module loader to use fake modules
    -- load opa:access module
    helpers = require("kong.plugins.opa.helpers")
  end)

  describe("filterHeaders function", function()
    it("three headers passed in and one matching wanted_headers", function()
      local headers = {}
      headers["unwantedheader1"] = "value1"
      headers["unwantedheader2"] = "value2"
      headers["wantedheader"] = "wantedvalue"

      local wanted_headers = {"wantedheader"}
      local res = helpers.filterHeaders(headers, wanted_headers)
      assert.are.same({["wantedheader"] = "wantedvalue"}, res)
    end)

    it("three headers passed in and one matching wanted_headers works regardless of upper/lower case", function()
      local headers = {}
      headers["unwantedheader1"] = "value1"
      headers["unwantedheader2"] = "value2"
      headers["wantedheader"] = "wantedValue"

      local wanted_headers = {"wantedHEADER"}
      local res = helpers.filterHeaders(headers, wanted_headers)
      assert.are.same({["wantedheader"] = "wantedValue"}, res)
    end)

    it("three headers passed in and zero matching wanted_headers", function()
      local headers = {}
      headers["unwantedheader1"] = "value1"
      headers["unwantedheader2"] = "value2"
      headers["unwantedheader4"] = "value3"

      local wanted_headers = {"wantedheader"}
      local res = helpers.filterHeaders(headers, wanted_headers)
      assert.are.same({}, res)
    end)

    it("three headers passed in and empty array for wanted_headers", function()
      local headers = {}
      headers["unwantedheader1"] = "value1"
      headers["unwantedheader2"] = "value2"
      headers["unwantedheader4"] = "value3"

      local wanted_headers = {}
      local res = helpers.filterHeaders(headers, wanted_headers)
      assert.are.same({}, res)
    end)

    it("empty headers table passed in and empty array for wanted_headers", function()
      local headers = {}
      headers["unwantedheader1"] = "value1"
      headers["unwantedheader2"] = "value2"
      headers["unwantedheader4"] = "value3"

      local wanted_headers = {}
      local res = helpers.filterHeaders(headers, wanted_headers)
      assert.are.same({}, res)
    end)
  end)

  describe("interp function", function()
    it("string tokens and parameters table match", function()
      local parameters_table = {}
      parameters_table["a"] = "Apple"
      parameters_table["b"] = "Banana"
      parameters_table["c"] = "Cherry"

      local tokenized_string = "fruit ${a}$ ${b} ${c} end"
      local res = helpers.interp(tokenized_string, parameters_table)
      assert.are.same("fruit Apple$ Banana Cherry end", res)
    end)

    it("string tokens and parameters table do not match", function()
      local parameters_table = {}
      parameters_table["a"] = "Apple"
      parameters_table["b"] = "Banana"
      parameters_table["c"] = "Cherry"

      local tokenized_string = "fruit ${x}$ ${y} ${z} end"
      local res = helpers.interp(tokenized_string, parameters_table)
      assert.are.same(tokenized_string, res)
    end)
  end)
end)
