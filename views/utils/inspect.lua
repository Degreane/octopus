-- Enhanced inspect function with options
function inspect(value, options)
    options = options or {}
    local depth = options.depth or 0
    local max_depth = options.max_depth or 10
    local visited = options.visited or {}
    local show_metatable = options.show_metatable or false
    local show_functions = options.show_functions or true
    local compact = options.compact or false
    
    local indent = compact and "" or string.rep("  ", depth)
    local newline = compact and "" or "\n"
    local space = compact and "" or " "
    
    -- Handle max depth
    if depth > max_depth then
        return indent .. "{ [max depth reached] }"
    end
    
    local value_type = type(value)
    
    if value_type == "table" then
        -- Check for circular references
        if visited[value] then
            return indent .. "{" .. space .. "[circular reference]" .. space .. "}"
        end
        visited[value] = true
        
        -- Check if table is empty
        local is_empty = true
        for _ in pairs(value) do
            is_empty = false
            break
        end
        
        if is_empty then
            visited[value] = nil
            return indent .. "{}"
        end
        
        local result = "{" .. newline
        local new_options = {
            depth = depth + 1,
            max_depth = max_depth,
            visited = visited,
            show_metatable = show_metatable,
            show_functions = show_functions,
            compact = compact
        }
        
        -- Sort keys for consistent output
        local keys = {}
        for k in pairs(value) do
            if show_functions or type(value[k]) ~= "function" then
                table.insert(keys, k)
            end
        end
        
        table.sort(keys, function(a, b)
            local ta, tb = type(a), type(b)
            if ta == tb then
                return tostring(a) < tostring(b)
            else
                return ta < tb
            end
        end)
        
        for i, k in ipairs(keys) do
            local v = value[k]
            local key_str
            
            if type(k) == "string" and k:match("^[%a_][%w_]*$") then
                key_str = k
            elseif type(k) == "string" then
                key_str = string.format('["%s"]', k:gsub('"', '\\"'))
            elseif type(k) == "number" then
                key_str = string.format("[%s]", k)
            else
                key_str = string.format("[%s]", inspect(k, new_options))
            end
            
            result = result .. indent .. "  " .. key_str .. space .. "=" .. space
            result = result .. inspect(v, new_options)
            
            if i < #keys then
                result = result .. ","
            end
            result = result .. newline
        end
        
        -- Show metatable if requested
        if show_metatable then
            local mt = getmetatable(value)
            if mt then
                result = result .. indent .. "  " .. "[metatable]" .. space .. "=" .. space
                result = result .. inspect(mt, new_options) .. newline
            end
        end
        
        visited[value] = nil -- Clean up for other branches
        return result .. indent .. "}"
        
    elseif value_type == "string" then
        return string.format('"%s"', value:gsub('"', '\\"'):gsub('\n', '\\n'):gsub('\t', '\\t'))
        
    elseif value_type == "function" then
        local info = debug.getinfo(value, "S")
        if info.what == "C" then
            return "[C function]"
        else
            return string.format("[function: %s:%d]", info.source, info.linedefined)
        end
        
    elseif value_type == "thread" then
        return string.format("[thread: %s]", tostring(value))
        
    elseif value_type == "userdata" then
        return string.format("[userdata: %s]", tostring(value))
        
    elseif value_type == "nil" then
        return "nil"
        
    elseif value_type == "boolean" then
        return value and "true" or "false"
        
    else
        return tostring(value)
    end
end

-- Convenience function for pretty printing
function pp(value, options)
    print(inspect(value, options))
end

-- Compact version
function inspect_compact(value)
    return inspect(value, { compact = true })
end


local Inspect = {
    inspect = inspect,
    pp = pp,
    inspect_compact = inspect_compact
}
return Inspect