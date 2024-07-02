/*******************************************************************************
 * @note Copyright (C) 2020 Shanghai Panchip Microelectronics Co., Ltd. All rights reserved.
 *
 * @file pan3028.h
 * @brief
 *
 * @history - V3.0, 2021-11-05
*******************************************************************************/
#ifndef __PAN_3028_H
#define __PAN_3028_H

 #include "pan3028_port.h"		

/* result */
#define OK                              0
#define FAIL                            1
/* 3028 mode define*/
#define PAN3028_MODE_DEEP_SLEEP         0
#define PAN3028_MODE_SLEEP              1
#define PAN3028_MODE_STB1               2
#define PAN3028_MODE_STB2               3
#define PAN3028_MODE_STB3               4
#define PAN3028_MODE_TX                 5
#define PAN3028_MODE_RX                 6
/* 3028 Tx mode */
#define PAN3028_TX_SINGLE               0
#define PAN3028_TX_CONTINOUS            1
/* 3028 Rx mode */
#define PAN3028_RX_SINGLE               0
#define PAN3028_RX_SINGLE_TIMEOUT       1
#define PAN3028_RX_CONTINOUS            2

/* System control register */
#define REG_SYS_CTL                     0x00
#define REG_FIFO_ACC_ADDR               0x01
/* 3V Logical area register */
#define REG_OP_MODE                     0X02

#define MODEM_MODE_NORMAL               0x01
#define MODEM_MODE_MULTI_SECTOR         0x02

#define freq_360000000                  (360000000)
#define freq_370000000                  (370000000)
#define freq_385000000                  (385000000)
#define freq_405000000                  (405000000)
#define freq_430000000                  (430000000)
#define freq_460000000                  (460000000)
#define freq_495000000                  (495000000)
#define freq_535000000                  (535000000)
#define freq_600000000                  (600000000)

#define freq_720000000                  (720000000)
#define freq_740000000                  (740000000)
#define freq_770000000                  (770000000)
#define freq_810000000                  (810000000)
#define freq_860000000                  (860000000)
#define freq_920000000                  (920000000)
#define freq_990000000                  (990000000)
#define freq_1070000000                 (1070000000)
#define freq_1200000000                 (1200000000)

#define CODE_RATE_45                     0x01
#define CODE_RATE_46                     0x02
#define CODE_RATE_47                     0x03
#define CODE_RATE_48                     0x04

#define SF_7                            7
#define SF_8                            8
#define SF_9                            9
#define SF_10                           10
#define SF_11                           11
#define SF_12                           12
    
#define BW_62_5K                        6
#define BW_125K                         7
#define BW_250K                         8
#define BW_500K                         9
    
#define CRC_OFF                         0
#define CRC_ON                          1
    
#define PLHD_IRQ_OFF                    0
#define PLHD_IRQ_ON                     1

#define PLHD_OFF                        0
#define PLHD_ON                         1

#define PLHD_LEN8                       0
#define PLHD_LEN16                      1

#define AGC_ON                          1
#define AGC_OFF                         0
	
#define LO_400M                         0
#define LO_800M                         1

#define DCDC_OFF                    	0
#define DCDC_ON                     	1

#define LDR_OFF                    		0
#define LDR_ON                     		1

#define LNA_GAIN_LOW                   	0
#define LNA_GAIN_HIGH                   1

#define CAD_DETECT_THRESHOLD_0A			0x0A
#define CAD_DETECT_THRESHOLD_10			0x10
#define CAD_DETECT_THRESHOLD_15			0x15
#define CAD_DETECT_THRESHOLD_20			0x20

#define REG_PAYLOAD_LEN                 0x0C

/*IRQ BIT MASK*/
#define REG_IRQ_RX_PLHD_DONE            0x10
#define REG_IRQ_RX_DONE                 0x8
#define REG_IRQ_CRC_ERR                 0x4
#define REG_IRQ_RX_TIMEOUT              0x2
#define REG_IRQ_TX_DONE                 0x1

enum REF_CLK_SEL {REF_CLK_32M,REF_CLK_16M};		
enum PAGE_SEL {PAGE0_SEL,PAGE1_SEL,PAGE2_SEL, PAGE3_SEL};

uint32_t PAN3028_rst(void);
uint32_t PAN3028_agc_enable(uint32_t state);
uint32_t PAN3028_agc_config(void);
uint32_t PAN3028_init(void);
uint32_t PAN3028_deepsleep_wakeup(void);
uint32_t PAN3028_deepsleep(void);
uint32_t PAN3028_sleep_wakeup(void);
uint32_t PAN3028_sleep(void);
uint32_t PAN3028_write_read_continue_regs(enum PAGE_SEL page,uint8_t addr,uint8_t *buffer,uint8_t len);
uint32_t PAN3028_set_freq(uint32_t freq);
uint32_t PAN3028_read_freq(void);
uint32_t PAN3028_calculate_tx_time(void);
uint32_t PAN3028_set_bw(uint32_t bw_val);
uint8_t PAN3028_get_bw(void);
uint32_t PAN3028_set_sf(uint32_t sf_val);
uint8_t PAN3028_get_sf(void);
uint32_t PAN3028_set_crc(uint32_t crc_val);
uint8_t PAN3028_get_crc(void);
uint8_t PAN3028_get_code_rate(void);
uint32_t PAN3028_set_code_rate(uint8_t code_rate);
uint32_t PAN3028_set_mode(uint8_t mode);
uint8_t PAN3028_get_mode(void);
uint32_t PAN3028_set_tx_mode(uint8_t mode);
uint32_t PAN3028_set_rx_mode(uint8_t mode);
uint32_t PAN3028_set_timeout(uint32_t timeout);
float PAN3028_get_snr(void);
float PAN3028_get_rssi(void);
uint32_t PAN3028_set_tx_power(uint8_t tx_power);
uint32_t PAN3028_get_tx_power(void);
uint32_t PAN3028_set_preamble(uint16_t reg);
uint32_t PAN3028_set_gpio_input(uint8_t gpio_pin);
uint32_t PAN3028_set_gpio_output(uint8_t gpio_pin);
uint32_t PAN3028_set_gpio_state(uint8_t gpio_pin, uint8_t state);
uint32_t PAN3028_cad_en(uint8_t threshold);
uint32_t PAN3028_cad_off(void);
uint32_t PAN3028_set_syncword(uint32_t sync);
uint8_t PAN3028_get_syncword(void);
uint32_t PAN3028_send_packet(uint8_t *buff, uint32_t len);
uint8_t PAN3028_recv_packet(uint8_t *buff);
uint32_t PAN3028_set_early_irq(uint32_t earlyirq_val);
uint8_t PAN3028_get_early_irq(void);
uint32_t PAN3028_set_plhd(uint8_t addr,uint8_t len);
uint8_t PAN3028_get_plhd(void);
uint32_t PAN3028_set_plhd_mask(uint32_t plhd_val);
uint8_t PAN3028_get_plhd_mask(void);
uint8_t PAN3028_recv_plhd8(uint8_t *buff);
uint8_t PAN3028_recv_plhd16(uint8_t *buff);
uint32_t PAN3028_plhd_receive(uint8_t *buf,uint8_t len);
uint32_t PAN3028_set_dcdc_mode(uint32_t dcdc_val);
uint32_t PAN3028_set_ldr(uint32_t mode);
uint32_t PAN3028_set_all_sf_preamble(uint32_t sf);
uint32_t PAN3028_set_all_sf_search(void);
uint32_t PAN3028_set_all_sf_search_off(void);
void PAN3028_irq_handler(void);
uint8_t PAN3028_get_irq(void);
uint8_t PAN3028_clr_irq(void);
uint32_t PAN3028_set_vco(uint8_t mode);
uint32_t PAN3028_set_lna_gain(uint8_t mode);
#endif
