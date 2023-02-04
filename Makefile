GOCMD=go
GOBUILD=${GOCMD} build -mod=mod
GOCLEAN=${GOCMD} clean

build: relationship

.PHONY: \
    relationbff relationship relationjob black blackjob relationbffprod relationshipprod relationjobprod blackprod blackjobprod

clean:
	${GOCLEAN}

relationship:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relationship github.com/zxq97/relation/app/relationship/service/cmd

relationjob:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_job github.com/zxq97/relation/app/relationship/job/cmd

relationbff:
	${GOBUILD} -o /Users/zongxingquan/goland/run/relation_bff github.com/zxq97/relation/app/relation/bff/cmd

black:
	${GOBUILD} -o /Users/zongxingquan/goland/run/black github.com/zxq97/relation/app/black/service/cmd

blackjob:
	${GOBUILD} -o /Users/zongxingquan/goland/run/black_job github.com/zxq97/relation/app/black/job/cmd

relationshipprod:
	${GOBUILD} -o /home/work/run/relationship github.com/zxq97/relation/app/relationship/service/cmd

relationjobprod:
	${GOBUILD} -o /home/work/run/relation_job github.com/zxq97/relation/app/relationship/job/cmd

relationbffprod:
	${GOBUILD} -o /home/work/run/relation_bff github.com/zxq97/relation/app/relation/bff/cmd

blackprod:
	${GOBUILD} -o /home/work/run/black github.com/zxq97/relation/app/black/service/cmd

blackjobprod:
	${GOBUILD} -o /home/work/run/black_job github.com/zxq97/relation/app/black/job/cmd
