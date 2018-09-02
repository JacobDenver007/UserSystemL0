-- 用合约来完成一个转账系统
local L0 = require("L0")

-- 合约创建时会被调用一次，之后就不会被调用
function L0Init(args)
    L0.PutState("created", "tostring(os.time())")
    print("Init OK")
    return true
end

-- 每次合约执行都调用
function L0Invoke(func, args)
    if("transfer" == func) then
        local cnt = #args
        print("cnt =",cnt)
        print("agr =",args)

        for i=0,#args do
            print("i=",i," key=",args[i])
            if (i%3==2)
            then
                local receiver= args[i-2]
                print("====>", receiver)
                local asset= args[i-1]
                print("====>", asset)
                local amount=tonumber(args[i])
                print("====>", amount)
                transfer(receiver, asset, amount)
            end
        end
    end
    print("Invoke OK")
    return true
end

-- 合约查询
function L0Query(args)
    return "Query ok"
end

function transfer(receiver, asset, amount)
    print("transfer", receiver, asset, amount)
    L0.Transfer(receiver, asset, amount)
end