#include "tpu_util.h"

int main(){
    int tpuNum = 0;
    getTpuNum(&tpuNum);
    printf("return %d \n",tpuNum);
	struct tpu_mem_info infoList[tpuNum];
    getTpuMemInfo(infoList,tpuNum);
    int i=0;
    for (i = 0; i < tpuNum; i++)
    {   
       
    }
    return 1;
};