#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdint.h>
#include <stdbool.h>
#include <sys/ioctl.h>
#include "bm_uapi.h"


struct tpu_mem_info{
    int dev_id;
	int chip_id;
	int chip_mode;  /*0---pcie; 1---soc*/
	int domain_bdf;
	int status;
	int card_index;
	int chip_index_of_card;

	int mem_used;
	int mem_total;
	int tpu_util;

	int board_temp;
	int chip_temp;
	int board_power;
	int tpu_power;
	int fan_speed;
};

int main(int argc, char **argv){
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

int getTpuNum(int *tpuNum){
      int arg;
      int fd;
      int ret;
      fd = open("/dev/bmdev-ctl", O_RDWR);
      if (fd < 0) {
        perror("open");
        exit(-2);
      }  
      ret = ioctl(fd, BMCTL_GET_DEV_CNT,&arg);
      if (ret) {
        perror("BMCTL_GET_DEV_CNT :");
        exit(-3);
      }
      *tpuNum = arg;
}




int getTpuMemInfo(struct tpu_mem_info infos[],int num){
    int fd;
    int ret;
    int arg;
    fd = open("/dev/bmdev-ctl", O_RDWR);
    if (fd < 0) {
        perror("open");
        exit(-2);
    }
    int i = 0;
    for (i = 0; i < num; i++)
    {   
        printf("获取第 %d 张设备信息",i);
        struct bm_smi_attr attr;
        struct tpu_mem_info info;
        attr.dev_id = i;
        ret = ioctl(fd, BMCTL_GET_SMI_ATTR,&attr);
        if (ret) {
          perror("BMCTL_GET_SMI_ATTR :");
          exit(-3);
        }
        info.dev_id = attr.dev_id;
        info.chip_id = attr.chip_id;
        info.chip_mode = attr.chip_mode;
        info.domain_bdf = attr.domain_bdf;
        info.card_index = attr.card_index;
        info.status = attr.status;
        info.chip_index_of_card = attr.chip_index_of_card;
        info.mem_used = attr.mem_used;
        info.mem_total = attr.mem_total;
        info.tpu_util = attr.tpu_util;
        info.board_temp = attr.board_temp;
        info.chip_temp = attr.chip_temp;
        info.board_power = attr.board_power;
        info.tpu_power = attr.tpu_power;
        info.fan_speed = attr.fan_speed;
        printf("dev_id %d \n",attr.dev_id);
        printf("mem_used %d \n",attr.mem_used);
        printf("mem_total %d \n",attr.mem_total);
        printf("card_index %d \n",attr.card_index);
        printf("board_temp %d \n",attr.board_temp);
        infos[i] = info;
    }
    return 1;
};
