# gServer_go
Server like gameServer build by go


1.WebServer
    1.GateServer register
    2.GateServer Player Online update
    3.5s check GateServer Online one time
    4.client get GateServer ip port info  keep same players on every GateServer
2.GateServer
    1.keep client online
    2.rec client msg
    3.send client msg to other function  server
    4.register function server
    5.unpack and check client msg
    6.send function msg to client
3.LoginServer
    1.client Login check 
    2.new client register
    3.role create delete
    4.update active client and clear long time unlogin client
    5.make cur login key and send key to GateServer
4.PublicServer
    1.deliver world chat msg
    2.deliver all server chat msg
    3.server notice
    4.e-mail
    5.friend sys 
