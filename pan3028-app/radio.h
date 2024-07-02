/*******************************************************************************
 * @note Copyright (C) 2020 Shanghai Panchip Microelectronics Co., Ltd. All rights reserved.
 *
 * @file radio.h
 * @brief
 *
 * @history - V3.0, 2021-11-05
*******************************************************************************/
#ifndef __RADIO_H_
#define __RADIO_H_

#include "pan3028_port.h"

//#define JAP_915
//#define FCC_915
//#define ETSI_433
//#define ETSI_868

#if defined(JAP_915)
#define DEFAULT_PWR            17
#define DEFAULT_FREQ           (915000000)
#define DEFAULT_SF             SF_7
#define DEFAULT_BW             BW_125K
#define DEFAULT_CR             CODE_RATE_48

#elif defined(FCC_915)
#define DEFAULT_PWR            16
#define DEFAULT_FREQ           (915000000)
#define DEFAULT_SF             SF_7
#define DEFAULT_BW             BW_125K
#define DEFAULT_CR             CODE_RATE_48

#elif defined(ETSI_433)
#define DEFAULT_PWR            13
#define DEFAULT_FREQ           (433000000)
#define DEFAULT_SF             SF_7
#define DEFAULT_BW             BW_125K
#define DEFAULT_CR             CODE_RATE_48

#elif defined(ETSI_868)
#define DEFAULT_PWR            12
#define DEFAULT_FREQ           (868000000)
#define DEFAULT_SF             SF_7
#define DEFAULT_BW             BW_125K
#define DEFAULT_CR             CODE_RATE_48

#else
/*Default parameter configuration*/
#define DEFAULT_PWR            26
#define DEFAULT_FREQ           (440250000)
#define DEFAULT_SF             SF_7
#define DEFAULT_BW             BW_250K
#define DEFAULT_CR             CODE_RATE_45
#endif


#define RADIO_FLAG_IDLE         0
#define RADIO_FLAG_TXDONE       1
#define RADIO_FLAG_RXDONE       2
#define RADIO_FLAG_RXTIMEOUT    3
#define RADIO_FLAG_RXERR        4
#define RADIO_FLAG_PLHDRXDONE   5

#define LEVEL_INACTIVE     		6
#define LEVEL_ACTIVE       		7

/* #define	GPIO_PIN_CAD			Pin02
#define	GPIO_PORT_CAD			PortB */

/* #define  CHECK_CAD()   			PORT_GetBit(GPIO_PORT_CAD, GPIO_PIN_CAD) */
#define  SET_TIMER_MS(time)		timer6_open_ms(time)	

#define MAC_EVT_TX_CAD_NONE			0x00
#define MAC_EVT_TX_CAD_TIMEOUT		0x01
#define MAC_EVT_TX_CAD_ACTIVE		0x02

struct RxDoneMsg
{
	uint8_t Payload[255];
	uint8_t PlhdPayload[16];
	uint16_t PlhdSize;
	uint16_t Size;
	float Rssi;
	float Snr;
};

typedef enum{
	RF_PARA_TYPE_FREQ,
	RF_PARA_TYPE_CR,
	RF_PARA_TYPE_BW,
	RF_PARA_TYPE_SF,
	RF_PARA_TYPE_TXPOWER,
	RF_PARA_TYPE_CRC,
}rf_para_type_t;

uint32_t rf_get_recv_flag(void);
void rf_set_recv_flag(int status);
uint32_t rf_get_transmit_flag(void);
void rf_set_transmit_flag(int status);
uint32_t rf_init(void);
uint32_t rf_deepsleep_wakeup(void);
uint32_t rf_deepsleep(void);
uint32_t rf_sleep_wakeup(void);
uint32_t rf_sleep(void);
uint32_t rf_get_tx_time(void);
uint32_t rf_set_mode(uint8_t mode);
uint8_t rf_get_mode(void);
uint32_t rf_set_tx_mode(uint8_t mode);
uint32_t rf_set_rx_mode(uint8_t mode);
uint32_t rf_set_rx_single_timeout(uint32_t timeout);
float rf_get_snr(void);
float rf_get_rssi(void);
uint32_t rf_set_preamble(uint16_t pream);
uint32_t rf_set_cad(uint8_t threshold);
uint32_t rf_set_cad_off(void);
uint32_t rf_set_syncword(uint8_t sync);
uint8_t rf_get_syncword(void);
void rf_irq_handler(void);
void rf_set_plhd_rx_on(uint8_t addr,uint8_t len);
void rf_set_plhd_rx_off(void);
uint32_t rf_receive(uint8_t *buf);
uint32_t rf_plhd_receive(uint8_t *buf,uint8_t len);
uint32_t rf_enter_continous_rx(void);
uint32_t rf_enter_single_timeout_rx(uint32_t timeout);
uint32_t rf_enter_single_rx(void);
uint32_t rf_single_tx_data(uint8_t *buf, uint8_t size, uint32_t *tx_time);
uint32_t rf_enter_continous_tx(void);
uint32_t rf_continous_tx_send_data(uint8_t *buf, uint8_t size);
uint32_t rf_set_agc(uint32_t state);
uint32_t rf_set_para(rf_para_type_t para_type, uint32_t para_val);
uint32_t rf_get_para(rf_para_type_t para_type, uint32_t *para_val);
void rf_set_default_para(void);
uint32_t rf_set_dcdc_mode(uint32_t dcdc_val);
uint32_t rf_set_ldr(uint32_t mode);
uint32_t rf_set_all_sf_preamble(uint32_t sf);
uint32_t rf_set_all_sf_search(void);
uint32_t rf_set_all_sf_search_off(void);
uint32_t rf_set_lna_gain(uint8_t mode);
void rf_irq_process(void);
uint32_t get_chirp_time(uint32_t bw,uint32_t sf);
/* uint32_t check_cad_rx_inactive(uint32_t one_chirp_time); */
/* uint32_t check_cad_tx_inactive(void); */
#endif

