# apiv1

call a v1 API, send body in stdin

        Usage of natip:
          -key string
                api v1 key
          -pass string
                api v1 password
          -method string
                http method
          -endpoint string
                url to call
          <STDIN>
                body contents on http request


## implementation

behavior: gets a v1 api session cookie via the passed credentials and invokes
the passed API call


## invocation

`echo '{"AccountAlias":"ECO","Name":"VA1ECOSRV01"}' | ./apiv1 -key XYZ -pass 'XYZ' -method GET -endpoint https://api.ctl.io/REST/Server/GetServer/JSON`
