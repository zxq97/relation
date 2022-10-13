GOCMD=go
GOBUILD=${GOCMD} build -mod=mod
GOCLEAN=${GOCMD} clean

build: relationsvc

.PHONY: \
    relationbff relationsvc relationtask relationadmin relationbffprod relationsvcprod relationtaskprod relationadminprod

clean:
	${GOCLEAN}

relationsvcprod:
	${GOBUILD} -o /home/work/run/relation_svc github.com/zxq97/relation/cmd/relationsvc

relationsvc:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_svc github.com/zxq97/relation/cmd/relationsvc

relationtaskprod:
	${GOBUILD} -o /home/work/run/relation_task github.com/zxq97/relation/cmd/relationtask

relationtask:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_task github.com/zxq97/relation/cmd/relationtask

relationbffprod:
	${GOBUILD} -o /home/work/run/relation_bff github.com/zxq97/relation/cmd/relationbff

relationbff:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_bff github.com/zxq97/relation/cmd/relationbff

relationadminprod:
	${GOBUILD} -o /home/work/run/relation_admin github.com/zxq97/relation/cmd/relationadmin

relationadmin:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_admin github.com/zxq97/relation/cmd/relationadmin
