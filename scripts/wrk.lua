wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

function request()
--     local n = math.random(1, 10000000) -- no cache
    local n = math.random(1, 200) -- cache
    local body = "username=" .. n .. "&password=" .. n
    wrk.body = body
    return wrk.format(wrk.method, "/login", wrk.headers, wrk.body)
end

--  wrk -t10 -c500 -d60s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t10 -c1000 -d60s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t20 -c1000 -d60s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t20 -c2000 -d20s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t20 -c2000 -d60s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t20 -c5000 -d60s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  wrk -t10 -c5000 -d20s -s ./scripts/wrk/wrk.lua http://localhost:8080/login
--  -t:test thread num; -c:connection num; -d:test duration;
--  go-torch -u http://localhost:8080/login -t 30
-- go tool pprof mybinary -http http://127.0.0.1:9999/debug/pprof/profile?seconds=30