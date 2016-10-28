# apiv1

call a v1 API, send body in stdin

        Usage of natip:
          -user string
                username
          -pass string
                password
          -method string
                http method
          -endpoint string
                url to call
          <STDIN>
                body contents on http request


## implementation

behavior: gets a v2 api session token via the passed credentials and invokes
the passed API call


## invocation

`echo '["VA1ECOSRV01"]' | ./apiv2 -user XYZ -pass 'XYZ' -method POST -endpoint https://api.ctl.io/v2/operations/ECO/servers/reboot`
