#include"stdio.h"

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

int getTpuNum(int *tpuNum);
int getTpuMemInfo(struct tpu_mem_info infos[],int num);
