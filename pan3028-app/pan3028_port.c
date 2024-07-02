/*******************************************************************************
 * @note Copyright (C) 2020 Shanghai Panchip Microelectronics Co., Ltd. All rights reserved.
 *
 * @file pan3028_port.c
 * @brief
 *
 * @history - V3.0, 2021-11-05
*******************************************************************************/
#ifdef KERNEL
#include <linux/delay.h>
#else
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/types.h>
#include <linux/spi/spidev.h>
#include <string.h>
#include <unistd.h>
#endif
#include "pan3028_port.h"
#include "spi-3028.h"

static uint8_t bits = 8; 
static uint32_t speed = 100000; 
static uint16_t delay = 10; 

extern uint8_t spi_tx_rx(uint8_t tx_data);

rf_port_t rf_port= 
{
	.antenna_init = rf_antenna_init,
	.tcxo_init = rf_tcxo_init,
	.set_tx = rf_antenna_tx,
	.set_rx = rf_antenna_rx,
	.antenna_close = rf_antenna_close,
	.tcxo_close = rf_tcxo_close,
	.spi_readwrite = spi_readwritebyte,
	.spi_cs_high = spi_cs_set_high,
	.spi_cs_low = spi_cs_set_low,
	.delayms = rf_delay_ms, 
	.delayus = rf_delay_us,
};

/**
 * @brief spi_readwritebyte
 * @param[in] <tx_data> spi readwritebyte value
 * @return result
 */
/* uint8_t spi_readwritebyte(uint8_t tx_data)
{
	int ret; 
	uint8_t tx[1] = {0};
	uint8_t rx[1] = {0}; //接收的数据数据

	return tx[0];
}
 */
/**
 * @brief spi_cs_set_high
 * @param[in] <none>
 * @return none
 */
void spi_cs_set_high(void)
{
/* 	PORT_SetBits(PortA, Pin04); */
}

/**
 * @brief spi_cs_set_low
 * @param[in] <none>
 * @return none
 */
void spi_cs_set_low(void)
{
/* 	PORT_ResetBits(PortA, Pin04); */
}

/**
 * @brief rf_delay_ms
 * @param[in] <time> ms
 * @return none
 */
void rf_delay_ms(uint32_t time)
{
	/* SysTick_Delay_10us(time*100); */
	// printf("sleep: %dus\n", time);
	zhg_usleep(time*1000);
}

/**
 * @brief rf_delay_us
 * @param[in] <time> us
 * @return none
 */
void rf_delay_us(uint32_t time)
{
	/* SysTick_Delay_10us(time/10+1); */
	// printf("sleep: %dus\n", time);
	zhg_usleep(time);
}

/**
 * @brief do PAN3028 TX/RX IO to initialize
 * @param[in] <none>
 * @return none
 */
void rf_antenna_init(void)
{
	PAN3028_set_gpio_output(MODULE_GPIO_RX);
	PAN3028_set_gpio_output(MODULE_GPIO_TX);

	PAN3028_set_gpio_state(MODULE_GPIO_RX, 0);
	PAN3028_set_gpio_state(MODULE_GPIO_TX, 0);    
}

/**
 * @brief do PAN3028 XTAL IO to initialize
 * @param[in] <none>
 * @return none
 */
void rf_tcxo_init(void)
{
	PAN3028_set_gpio_output(MODULE_GPIO_TCXO);
	PAN3028_set_gpio_state(MODULE_GPIO_TCXO, 1);
}

/**
 * @brief close PAN3028 XTAL IO 
 * @param[in] <none>
 * @return none
 */
void rf_tcxo_close(void)
{
	PAN3028_set_gpio_output(MODULE_GPIO_TCXO);
	PAN3028_set_gpio_state(MODULE_GPIO_TCXO, 0);
}
/**
 * @brief change PAN3028 IO to rx
 * @param[in] <none>
 * @return none
 */
void rf_antenna_rx(void)
{ 
	PAN3028_set_gpio_state(MODULE_GPIO_TX, 0);     
	PAN3028_set_gpio_state(MODULE_GPIO_RX, 1);
}

/**
 * @brief change PAN3028 IO to tx
 * @param[in] <none>
 * @return none
 */
void rf_antenna_tx(void)
{
	PAN3028_set_gpio_state(MODULE_GPIO_RX, 0);      
	PAN3028_set_gpio_state(MODULE_GPIO_TX, 1);
}

/**
 * @brief change PAN3028 IO to close
 * @param[in] <none>
 * @return none
 */
void rf_antenna_close(void)
{
	PAN3028_set_gpio_state(MODULE_GPIO_TX, 0);
	PAN3028_set_gpio_state(MODULE_GPIO_RX, 0);
}

