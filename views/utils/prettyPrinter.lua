local PrettyPrinter = {}

-- Main function to determine type and print accordingly
function PrettyPrinter.print(input)
    local inputType = type(input)
    
    if inputType == "string" then
        print(input)
    elseif inputType == "table" then
        PrettyPrinter.prettyPrintTable(input, 0)
    elseif inputType == "number" then
        print(tostring(input))
    elseif inputType == "boolean" then
        print(tostring(input))
    elseif inputType == "nil" then
        print("nil")
    elseif inputType == "function" then
        print("function: " .. tostring(input))
    else
        print(inputType .. ": " .. tostring(input))
    end
end

-- Pretty print function for tables with indentation
function PrettyPrinter.prettyPrintTable(tbl, indent)
    indent = indent or 0
    local indentStr = string.rep("  ", indent)
    
    -- Check if table is empty
    if next(tbl) == nil then
        print(indentStr .. "{}")
        return
    end
    
    print(indentStr .. "{")
    
    -- Separate numeric and string keys for better formatting
    local numericKeys = {}
    local stringKeys = {}
    local otherKeys = {}
    
    for key, _ in pairs(tbl) do
        if type(key) == "number" then
            table.insert(numericKeys, key)
        elseif type(key) == "string" then
            table.insert(stringKeys, key)
        else
            table.insert(otherKeys, key)
        end
    end
    
    -- Sort keys for consistent output
    table.sort(numericKeys)
    table.sort(stringKeys)
    
    -- Print numeric keys first
    for _, key in ipairs(numericKeys) do
        PrettyPrinter.printKeyValue(key, tbl[key], indent + 1)
    end
    
    -- Then string keys
    for _, key in ipairs(stringKeys) do
        PrettyPrinter.printKeyValue(key, tbl[key], indent + 1)
    end
    
    -- Finally other key types
    for _, key in ipairs(otherKeys) do
        PrettyPrinter.printKeyValue(key, tbl[key], indent + 1)
    end
    
    print(indentStr .. "}")
end

-- Helper function to print key-value pairs
function PrettyPrinter.printKeyValue(key, value, indent)
    local indentStr = string.rep("  ", indent)
    local keyStr = PrettyPrinter.formatKey(key)
    local valueType = type(value)
    
    if valueType == "table" then
        print(indentStr .. keyStr .. " = ")
        PrettyPrinter.prettyPrintTable(value, indent + 1)
    elseif valueType == "string" then
        print(indentStr .. keyStr .. ' = "' .. value .. '"')
    elseif valueType == "number" or valueType == "boolean" then
        print(indentStr .. keyStr .. " = " .. tostring(value))
    elseif valueType == "nil" then
        print(indentStr .. keyStr .. " = nil")
    elseif valueType == "function" then
        print(indentStr .. keyStr .. " = function: " .. tostring(value))
    else
        print(indentStr .. keyStr .. " = " .. valueType .. ": " .. tostring(value))
    end
end

-- Helper function to format keys properly
function PrettyPrinter.formatKey(key)
    if type(key) == "string" then
        -- Check if key is a valid Lua identifier
        if string.match(key, "^[%a_][%w_]*$") then
            return key
        else
            return '["' .. key .. '"]'
        end
    elseif type(key) == "number" then
        return "[" .. tostring(key) .. "]"
    else
        return "[" .. tostring(key) .. "]"
    end
end

-- Alternative compact version for single line output
function PrettyPrinter.printCompact(input)
    local inputType = type(input)
    
    if inputType == "string" then
        print('"' .. input .. '"')
    elseif inputType == "table" then
        print(PrettyPrinter.tableToString(input))
    else
        print(tostring(input))
    end
end

-- Helper function to convert table to compact string
function PrettyPrinter.tableToString(tbl)
    if next(tbl) == nil then
        return "{}"
    end
    
    local parts = {}
    local isArray = PrettyPrinter.isArray(tbl)
    
    if isArray then
        for i = 1, #tbl do
            local value = tbl[i]
            if type(value) == "table" then
                table.insert(parts, PrettyPrinter.tableToString(value))
            elseif type(value) == "string" then
                table.insert(parts, '"' .. value .. '"')
            else
                table.insert(parts, tostring(value))
            end
        end
        return "{" .. table.concat(parts, ", ") .. "}"
    else
        for key, value in pairs(tbl) do
            local keyStr = PrettyPrinter.formatKey(key)
            local valueStr
            
            if type(value) == "table" then
                valueStr = PrettyPrinter.tableToString(value)
            elseif type(value) == "string" then
                valueStr = '"' .. value .. '"'
            else
                valueStr = tostring(value)
            end
            
            table.insert(parts, keyStr .. " = " .. valueStr)
        end
        return "{" .. table.concat(parts, ", ") .. "}"
    end
end

-- Helper function to check if table is an array
function PrettyPrinter.isArray(tbl)
    local count = 0
    for _ in pairs(tbl) do
        count = count + 1
    end
    return count == #tbl
end

return PrettyPrinter
