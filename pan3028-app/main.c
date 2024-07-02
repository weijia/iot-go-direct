#include <stdint.h>
#include <stdio.h>
#include <sys/ioctl.h>
#include "pan3028.h"

#define TX_LEN  10
uint8_t tx_test_buf[TX_LEN] = {0, 1, 2, 3, 4, 5, 6, 7, 8, 9};

int main(void){
    int ret = 0;
	int cnt = 0;
	ret = rf_init();
	if(ret != OK)
	{
		printf(" RF Init Fail");
		while(1);
	}	
	rf_set_syncword(0x12);
	rf_set_default_para();
	rf_enter_continous_tx();
	if(rf_continous_tx_send_data(tx_test_buf, TX_LEN) != OK)		
	{
		printf("tx fail \r\n");
	}
	else
	{
		cnt ++;
		printf("Tx cnt %d\r\n", cnt );
	}
				
	while (1)
	{
		rf_irq_process();
		if(rf_get_transmit_flag() == RADIO_FLAG_TXDONE)
		{
/* 			BSP_LED_Toggle(); */
			rf_set_transmit_flag(RADIO_FLAG_IDLE);  			
/* 			SysTick_Delay(1000); */
			if(rf_continous_tx_send_data(tx_test_buf, TX_LEN) != OK)		
			{
				printf("tx fail \r\n");
			}
			else
			{
				cnt ++;
				printf("Tx cnt %d\r\n", cnt );
			}
		}
	}
}