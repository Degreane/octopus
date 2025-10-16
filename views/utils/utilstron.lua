local pp = require("views.utils.prettyPrinter")

-- Execute node tron.js command to estimate TRC20 transaction fees
local function estimateTxTrc20Fees(contractAddress, walletAddress, amount, pkey)
    if not contractAddress or contractAddress == "" then
        return nil, "Contract address is required"
    end
    if not walletAddress or walletAddress == "" then
        return nil, "Wallet address is required"
    end
    if not amount then
        amount = 1000000.0
    end
    if not pkey or pkey == "" then
        return nil, "Private key is required"
    end

    -- Format the command string with all arguments
    local command = string.format(
        "node tron.js --estimateTxTrc20Fees --contractAddress=%s --account=%s --amount=%f --pkey=%s", contractAddress,
        walletAddress, amount, pkey)
    print(string.format("Executing command-1: %s", command))
    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        return nil, "Failed to execute command"
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()

    if not success then
        return nil, string.format("Command failed with exit code: %s", exit_code or "unknown")
    end

    -- Clean up the result - extract JSON if present
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim whitespace

    -- Remove command path if it appears in output
    result = result:gsub("^/usr/bin/node", "")
    result = result:gsub("^/usr/local/bin/node", "")
    result = result:gsub("^node", "")
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim again

    -- Extract JSON from the output
    local json_match = result:match('({.*})')
    if json_match then
        return json_match, nil
    end

    return result, nil
end
local function getAccount(walletAddress)
    if not walletAddress or walletAddress == "" then
        return nil, "Wallet address is required"
    end

    -- Log the operation
    print(string.format("Getting account information for wallet: %s", walletAddress))

    -- Format the command string
    local command = string.format("node tron.js --getAccount %s", walletAddress)
    print(string.format("Executing command: %s", command))

    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        local error_msg = "Failed to execute command"
        print("Error: " .. error_msg)
        return nil, error_msg
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()

    if not success then
        local error_msg = string.format("Command failed with exit code: %s", exit_code or "unknown")
        print("Error: " .. error_msg)
        return nil, error_msg
    end

    -- Clean up the result - extract JSON if present
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim whitespace

    -- Remove command path if it appears in output
    result = result:gsub("^/usr/bin/node", "")
    result = result:gsub("^/usr/local/bin/node", "")
    result = result:gsub("^node", "")
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim again

    -- Extract JSON from the output
    local json_match = result:match('({.*})')
    if json_match then
        print(string.format("Account info result: %s", json_match))
        return json_match, nil
    end

    print(string.format("Raw result: %s", result))
    return result, nil
end

local function getAccountBalanceWithLogging(walletAddress)
    if not walletAddress or walletAddress == "" then
        return nil, "Wallet address is required"
    end

    -- Log the operation
    print(string.format("Getting account balance for wallet: %s", walletAddress))

    -- Format the command string
    local command = string.format("node tron.js --getAccountBalance=%s", walletAddress)
    -- local command = string.format("pwd")
    print(string.format("Executing command: %s", command))

    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        local error_msg = "Failed to execute command"
        print("Error: " .. error_msg)
        return nil, error_msg
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()

    if not success then
        local error_msg = string.format("Command failed with exit code: %s", exit_code or "unknown")
        print("Error: " .. error_msg)
        return nil, error_msg
    end

    -- Trim whitespace from result
    result = result:gsub("^%s*(.-)%s*$", "%1")

    print(string.format("Command result: %s", result))
    return result, nil
end

local function getTrxPrice()
    --[[
        url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&per_page=250&sparkline=true&price_change_percentage=1h%2C%2024h%2C%207d%2C%2014d%2C%20200d%2C%201y&precision=full"
        curl --request GET \
            --url 'https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=tron&names=Tron&symbols=trx&include_tokens=top&x_cg_demo_api_key=CG-QHHMB8UJdAQrbFsdLhVLJ7D3' \
            --header 'accept: application/json'
        req.Header.Add("accept", "application/json")
            req.Header.Add("x-cg-demo-api-key", "CG-QHHMB8UJdAQrbFsdLhVLJ7D3")
    ]]

    local trx_price = 0.9 -- trx_price is the price in US Dollars for TRX
    local CoinGecko = eocto.makeRequest("GET",
        'https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=tron&names=Tron&symbols=trx&include_tokens=top&x_cg_demo_api_key=CG-QHHMB8UJdAQrbFsdLhVLJ7D3',
        {"accept: application/json", "x-cg-demo-api-key: CG-QHHMB8UJdAQrbFsdLhVLJ7D3"})
    if CoinGecko["status"] == 200 then
        trx_price = eocto.decodeJSON(CoinGecko["body"])[1].current_price
    end
    return trx_price
end
-- Send TRX to another address
-- @param receiverAddress: The TRON address to receive the TRX
-- @param senderPrivateKey: The private key of the sender
-- @param amountTRX: The amount to send in TRX (will be converted to SUN)
local function sendTRX(receiverAddress, senderPrivateKey, amountTRX)
    if not receiverAddress or receiverAddress == "" then
        return nil, "Receiver address is required"
    end
    if not senderPrivateKey or senderPrivateKey == "" then
        return nil, "Private key is required"
    end
    if not amountTRX or amountTRX <= 0 then
        return nil, "Amount must be greater than 0"
    end

    -- Convert TRX to SUN (1 TRX = 1,000,000 SUN)
    local amountSUN = amountTRX * 1000000

    -- Format the command string with all arguments
    local command = string.format("node tron.js --sendTRX --account=%s --pkey=%s --amount=%.0f", receiverAddress,
        senderPrivateKey, amountSUN)

    print(string.format("Executing sendTRX command: %s", command))

    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        return nil, "Failed to execute command"
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()

    if not success then
        return nil, string.format("Command failed with exit code: %s", exit_code or "unknown")
    end

    -- Clean up the result - extract JSON if present
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim whitespace

    -- Remove command path if it appears in output
    result = result:gsub("^/usr/bin/node", "")
    result = result:gsub("^/usr/local/bin/node", "")
    result = result:gsub("^node", "")
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim again

    -- Extract JSON from the output
    local json_match = result:match('({.*})')
    if json_match then
        print(string.format("SendTRX result: %s", json_match))
        return json_match, nil
    end

    print(string.format("SendTRX raw result: %s", result))
    return result, nil
end

-- Send TRC20 tokens to another address
-- @param contractAddress: The TRC20 contract address (e.g., USDT contract)
-- @param recipientAddress: The TRON address to receive the tokens
-- @param senderPrivateKey: The private key of the sender
-- @param amount: The amount to send in token units
local function sendTRC20(contractAddress, recipientAddress, senderPrivateKey, amount)
    if not contractAddress or contractAddress == "" then
        return nil, "Contract address is required"
    end
    if not recipientAddress or recipientAddress == "" then
        return nil, "Recipient address is required"
    end
    if not senderPrivateKey or senderPrivateKey == "" then
        return nil, "Private key is required"
    end
    if not amount or amount <= 0 then
        return nil, "Amount must be greater than 0"
    end

    -- Convert token amount to smallest unit (similar to TRX to SUN conversion)
    -- Most TRC20 tokens use 6 decimal places (like USDT), so multiply by 1,000,000
    local amountInSmallestUnit = amount * 1000000
    -- if true then
    --     return eocto.encodeJSON({
    --         txid = 'acd950206defff5808e11a1d315d2bda2a26a6c2642bffc341c07e5d540d308d',
    --         result = true
    --     })
    -- end
    -- Format the command string with all arguments
    local command = string.format("node tron.js --sendTRC20 --contractAddress=%s --account=%s --amount=%.f --pkey=%s",
        contractAddress, recipientAddress, amountInSmallestUnit, senderPrivateKey)

    print(string.format("Executing sendTRC20 command: %s", command))

    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        return nil, "Failed to execute command"
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()

    if not success then
        return nil, string.format("Command failed with exit code: %s", exit_code or "unknown")
    end

    -- Clean up the result - extract JSON if present
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim whitespace

    -- Remove command path if it appears in output
    result = result:gsub("^/usr/bin/node", "")
    result = result:gsub("^/usr/local/bin/node", "")
    result = result:gsub("^node", "")
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim again

    -- Extract JSON from the output
    local json_match = result:match('({.*})')
    if json_match then
        print(string.format("SendTRC20 result: %s", json_match))
        return json_match, nil
    end

    print(string.format("SendTRC20 raw result: %s", result))
    return result, nil
end

-- Get transaction information by transaction ID
-- @param txid: The transaction ID to query
-- @return result, error
local function getTransactionInfo(txid)
    if not txid or txid == "" then
        return nil, "Transaction ID is required"
    end
    -- Strip and clean the txid of hidden/unseen characters
    txid = string.gsub(txid, "%s", "") -- Remove all whitespace (spaces, tabs, newlines)
    txid = string.gsub(txid, "%c", "") -- Remove all control characters
    txid = string.gsub(txid, "[%z\1-\31]", "") -- Remove null and control characters (0-31)
    txid = string.gsub(txid, "[^\32-\126]", "") -- Keep only printable ASCII characters (32-126)
    txid = string.gsub(txid, "^%s*(.-)%s*$", "%1") -- Trim leading/trailing whitespace

    -- Additional cleaning for common hidden characters
    txid = string.gsub(txid, "\r", "") -- Remove carriage returns
    txid = string.gsub(txid, "\n", "") -- Remove newlines
    txid = string.gsub(txid, "\t", "") -- Remove tabs
    txid = string.gsub(txid, "\0", "") -- Remove null bytes

    -- Validate txid length (TRON transaction IDs are typically 64 characters)
    if string.len(txid) ~= 64 then
        return nil, "Invalid transaction ID length. Expected 64 characters, got " .. string.len(txid)
    end

    -- Validate txid contains only hexadecimal characters
    if not string.match(txid, "^[a-fA-F0-9]+$") then
        return nil, "Invalid transaction ID format. Must contain only hexadecimal characters"
    end
    -- Log the operation
    print(string.format("Getting transaction info for txid: %s", txid))

    -- Format the command string
    local command = string.format("node tron.js --getTransactionInfo='%s' ", txid)
    print(string.format("Executing command: %s", command))

    -- Execute the command and capture output
    local handle = io.popen(command)
    if not handle then
        local error_msg = "Failed to execute command"
        print("Error: " .. error_msg)
        return nil, error_msg
    end

    local result = handle:read("*a") -- Read all output
    local success, exit_type, exit_code = handle:close()
    -- print("Success")
    -- pp.print(success)
    -- pp.print(exit_type)
    -- pp.print(exit_code)
    -- print("----------------------^^^^^^^^^^^^^^^^^^^^--------------------")
    if not success then
        local error_msg = string.format("Command failed with exit code: %s", exit_code or "unknown")
        print("Error: " .. error_msg)
        return nil, error_msg
    end
    eocto.debug("warning", "getTransactionInfo result: " .. result)
    -- Clean up the result - extract JSON if present
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim whitespace

    -- Remove command path if it appears in output
    result = result:gsub("^/usr/bin/node", "")
    result = result:gsub("^/usr/local/bin/node", "")
    result = result:gsub("^node", "")
    result = result:gsub("^%s*(.-)%s*$", "%1") -- Trim again

    -- Extract JSON from the output
    local json_match = result:match('({.*})')

    if json_match then
        print(string.format("Transaction info result: %s", json_match))

        -- Parse the JSON response and return the parsed object
        -- local response = eocto.decodeJSON(json_match)
        -- if response then
        --     -- print("Transaction info response:")
        --     -- pp.print(response)
        --     return response, nil
        -- else
            -- print("Failed to parse JSON response")
        return json_match, nil -- Return raw JSON if parsing fails
        -- end
    end

    print(string.format("Raw result: %s", result))
    return result, nil
end

local tronutils = {
    getAccountBalanceWithLogging = getAccountBalanceWithLogging,
    estimateTxTrc20Fees = estimateTxTrc20Fees,
    getAccount = getAccount,
    getTrxPrice = getTrxPrice,
    sendTRX = sendTRX,
    sendTRC20 = sendTRC20,
    getTransactionInfo = getTransactionInfo
}

return tronutils
