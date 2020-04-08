local _M = {}

function _M.load_jwt(self, jwt_str, secret)
    return {
        payload = {
            sub = "user_id",
            role = { "user_role1", "user_role2" },
        }
    }
end

return _M