---
--- Created by fbanna.
--- DateTime: 8/15/25 6:12 PM
---
local pp = require('views.utils.prettyPrinter')

local daysLimits = {
    [360] = {
        ["color"] = "text-[#008CC9]",
        ["title"] = "Your Membership Is Active",
        ["ccolor"] = "#008CC9"
    },
    [180] = {
        ["color"] = "text-[#FFC83D]",
        ["title"] = "Your Renewal Date Is Approaching",
        ["ccolor"] = "#FFC83D"
    },
    [90] = {
        ["color"] = "text-[#FF6208]",
        ["title"] = "You’re Running Out of Time",
        ["ccolor"] = "#FF6208"
    },
    [1] = {
        ["color"] = "text-[#FB1D00]",
        ["title"] = "Your Membership Is About to Expire",
        ["ccolor"] = "#FB1D00"
    }
}
eid = tostring(eocto.getPathParam("id")) or nil
local pvState = {}
local daysleft = 0
if eid ~= nil then
    -- if the eid is not nil then we should make a request to get the pvState
    local pvRequest = eocto.makeRequest(
            "GET",
            string.format("https://egety.me/api/pl3nv!d4/pv_5t4t3/%s", eid),
            {
                ["Content-Type"] = "application/json"
            }
    )
    if pvRequest and pvRequest["status"] == 200 then
        local pvResponse = eocto.decodeJSON(pvRequest["body"])
        if pvResponse["success"] == true then
            daysleft = pvResponse["data"]["remainingDays"]
            pvState = pvResponse["data"]
            pp.print(pvState)
        end
    end
end

local daysLimit = {}

if daysleft <= 0 then
    daysleft = 0
elseif daysleft > 360 then
    daysLimit = daysLimits[360]
elseif daysleft > 180 then
    daysLimit = daysLimits[180]
elseif daysleft > 90 then
    daysLimit = daysLimits[90]
else
    daysLimit = daysLimits[1]
end
eocto.setLocal("daysLeft", daysleft)
eocto.setLocal("daysLimit", daysLimit)
eocto.setLocal("pvState", pvState)

