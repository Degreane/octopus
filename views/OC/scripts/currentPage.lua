local Inspect = require('views.utils.inspect')
local path=eocto.getPath()
print("path: " .. path)
local paths={
    ["/OC/"]="home",
    ["/OC/lua"]="lua",
    ["/OC/yaml"]="yaml",
    ["/OC/helpers"]="helpers",
    ["/about"]="about",
    ["/contact"]="contact",
    ["/services"]="services",
    ["/portfolio"]="portfolio",
    ["/blog"]="blog",
    ["/blog/post"]="post",
}
if path ~= nil then 
    eocto.setLocal("currentPage",paths[path])
end