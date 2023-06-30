package.path=package.path..";C:\\roles\\shared\\go\\src\\github.com\\RiftTeam\\tic-80-bundler\\lua-tests\\example4\\?.lua"

L={}
require("./file1")(L)
require("./file2")(L)

function TIC()
    cls()
    local val=L.calculate(23,6)
    L.print(val)
end

-- <PALETTE>
-- 000:1a1c2c5d275db13e53ef7d57ffcd75a7f07038b76425717929366f3b5dc941a6f673eff7f4f4f494b0c2566c86333c57
-- </PALETTE>
