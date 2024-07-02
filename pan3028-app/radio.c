/*******************************************************************************
 * @note Copyright (C) 2020 Shanghai Panchip Microelectronics Co., Ltd. All rights reserved.
 *
 * @file radio.c
 * @brief
 *
 * @history - V3.0, 2021-11-05
*******************************************************************************/
#ifdef KERNEL
#include <linux/init.h>
#include <linux/module.h>
#include <linux/ioctl.h>
#include <linux/fs.h>
#include <linux/device.h>
#include <linux/err.h>
#include <linux/list.h>
#include <linux/errno.h>
#include <linux/mutex.h>
#include <linux/slab.h>
#include <linux/compat.h>
#include <linux/of.h>
#include <linux/of_device.h>
#endif
#include "pan3028_port.h"
#include "spi-3028.h"

extern bool pan3028_irq_trigged_flag;
/*
 * flag that indicate if a new packet is received.
*/
static int packet_received = RADIO_FLAG_IDLE;

/*
 * flag that indicate if transmision is finished.
*/
static int packet_transmit = RADIO_FLAG_IDLE;

struct RxDoneMsg RxDoneParams;

uint32_t cad_tx_timeout_flag = MAC_EVT_TX_CAD_NONE;
/**
 * @brief get receive flag 
 * @param[in] <none>
 * @return receive state
 */
uint32_t rf_get_recv_flag(void)
{
	return packet_received;
}

/**
 * @brief set receive flag 
 * @param[in] <status> receive flag state to set
 * @return none
 */
void rf_set_recv_flag(int status)
{
	packet_received = status;
}

/**
 * @brief get transmit flag 
 * @param[in] <none>
 * @return reansmit state
 */
uint32_t rf_get_transmit_flag(void)
{
	return packet_transmit;
}

/**
 * @brief set transmit flag 
 * @param[in] <status> transmit flag state to set
 * @return none
 */
void rf_set_transmit_flag(int status)
{
	packet_transmit = status;
}

/**
 * @brief do basic configuration to initialize
 * @param[in] <none>
 * @return result
 */
uint32_t rf_init(void)
{
	pr_info("rf_init\n");
	if(PAN3028_deepsleep_wakeup() != OK)
	{
		return FAIL;
	} 
	pr_info("PAN3028_init\n");
	if(PAN3028_init() != OK)
	{
		return FAIL;
	} 

	if(rf_set_agc(AGC_ON) != OK)
	{
		return FAIL;
	} 

	rf_port.antenna_init();
	pr_info("rf_init OK\n");
	return OK;    
}

/**
 * @brief change PAN3028 mode from deep sleep to wakeup(STB3)
 * @param[in] <none>
 * @return result
 */
uint32_t rf_deepsleep_wakeup(void)
{
	if(PAN3028_deepsleep_wakeup() != OK)
	{
		return FAIL;
	} 

	if(PAN3028_init() != OK)
	{
		return FAIL;
	} 

	if(rf_set_agc(AGC_ON) != OK)
	{
		return FAIL;
	} 

	rf_port.antenna_init();

	return OK; 
}

/**
 * @brief change PAN3028 mode from sleep to wakeup(STB3) 
 * @param[in] <none>
 * @return result
 */
uint32_t rf_sleep_wakeup(void)
{
	if(PAN3028_sleep_wakeup() != OK)
	{
		return FAIL;
	}     
	rf_port.antenna_init();
	return OK;
}

/**
 * @brief change PAN3028 mode from standby3(STB3) to deep sleep, PAN3028 should set DCDC_OFF before enter deepsleep
 * @param[in] <none>
 * @return result
 */
uint32_t rf_deepsleep(void)
{
	rf_port.antenna_close();
	return PAN3028_deepsleep();
}

/**
 * @brief change PAN3028 mode from standby3(STB3) to deep sleep, PAN3028 should set DCDC_OFF before enter sleep
 * @param[in] <none>
 * @return result
 */
uint32_t rf_sleep(void)
{
	rf_port.antenna_close();
	return PAN3028_sleep();
}
	
/**
 * @brief calculate tx time
 * @param[in] <none>
 * @return tx time(ms) 
 */
uint32_t rf_get_tx_time(void)
{
#ifndef KERNEL
	return PAN3028_calculate_tx_time();
#else
	return 1;
#endif
}

/**
 * @brief set rf mode
 * @param[in] <mode>    
 *			  PAN3028_MODE_DEEP_SLEEP / PAN3028_MODE_SLEEP
 *            PAN3028_MODE_STB1 / PAN3028_MODE_STB2
 *            PAN3028_MODE_STB3 / PAN3028_MODE_TX / PAN3028_MODE_RX
 * @return result
 */
uint32_t rf_set_mode(uint8_t mode)
{
	return PAN3028_set_mode(mode);
}

/**
 * @brief get rf mode
 * @param[in] <none>
 * @return mode 
 *		   PAN3028_MODE_DEEP_SLEEP / PAN3028_MODE_SLEEP
 *         PAN3028_MODE_STB1 / PAN3028_MODE_STB2
 *         PAN3028_MODE_STB3 / PAN3028_MODE_TX / PAN3028_MODE_RX
 */
uint8_t rf_get_mode(void)
{
	return PAN3028_get_mode();
}

/**
 * @brief set rf Tx mode
 * @param[in] <mode> 
 *			  PAN3028_TX_SINGLE/PAN3028_TX_CONTINOUS
 * @return result
 */
uint32_t rf_set_tx_mode(uint8_t mode)
{
	return PAN3028_set_tx_mode(mode);
}

/**
 * @brief set rf Rx mode
 * @param[in] <mode> 
 *			  PAN3028_RX_SINGLE/PAN3028_RX_SINGLE_TIMEOUT/PAN3028_RX_CONTINOUS
 * @return result
 */
uint32_t rf_set_rx_mode(uint8_t mode)
{
	return PAN3028_set_rx_mode(mode);
}

/**
 * @brief set timeout for Rx. It is useful in PAN3028_RX_SINGLE_TIMEOUT mode
 * @param[in] <timeout> rx single timeout time(in ms)
 * @return result
 */
uint32_t rf_set_rx_single_timeout(uint32_t timeout)
{
	return PAN3028_set_timeout(timeout);
}

/**
 * @brief get snr value
 * @param[in] <none> 
 * @return snr
 */
float rf_get_snr(void)
{
#ifdef KERNEL
	return PAN3028_get_snr();
#else
	return 1.0;
#endif
}

/**
 * @brief get rssi value
 * @param[in] <none> 
 * @return rssi
 */
float rf_get_rssi(void)
{
#ifdef KERNEL
	return PAN3028_get_rssi();
#else
	return 1.0;
#endif
}

/**
 * @brief set preamble 
 * @param[in] <reg> preamble
 * @return result
 */
uint32_t rf_set_preamble(uint16_t pream)
{
	return PAN3028_set_preamble(pream);
}

/**
 * @brief CAD function enable
 * @param[in] <threshold> 
			  CAD_DETECT_THRESHOLD_0A / CAD_DETECT_THRESHOLD_10 / CAD_DETECT_THRESHOLD_15 / CAD_DETECT_THRESHOLD_20
 * @return  result
 */
uint32_t rf_set_cad(uint8_t threshold)
{
	return PAN3028_cad_en(threshold);
}

/**
 * @brief CAD function disable
 * @param[in] <none> 
 * @return  result
 */
uint32_t rf_set_cad_off(void)
{
	return PAN3028_cad_off();
}

/**
 * @brief set rf syncword
 * @param[in] <sync> syncword
 * @return result
 */
uint32_t rf_set_syncword(uint8_t sync)
{
	return PAN3028_set_syncword(sync);
}

/**
 * @brief read rf syncword
 * @param[in] <none>   
 * @return syncword
 */
uint8_t rf_get_syncword(void)
{
	return PAN3028_get_syncword();
}

/**
 * @brief RF IRQ server routine, it should be call at ISR of IRQ pin
 * @param[in] <none>
 * @return result
 */
void rf_irq_handler(void)
{
	PAN3028_irq_handler();
}

/**
 * @brief set rf plhd mode on , rf will use early interruption
 * @param[in] <addr> PLHD start addr,Range:0..7f
		      <len> PLHD len
			  PLHD_LEN8 / PLHD_LEN16
 * @return result
 */
void rf_set_plhd_rx_on(uint8_t addr,uint8_t len)
{
	PAN3028_set_early_irq(PLHD_IRQ_ON);
	PAN3028_set_plhd(addr,len);
	PAN3028_set_plhd_mask(PLHD_ON);
}

/**
 * @brief set rf plhd mode off
 * @param[in] <none>
 * @return result
 */
void rf_set_plhd_rx_off(void)
{
	PAN3028_set_early_irq(PLHD_IRQ_OFF);
	PAN3028_set_plhd_mask(PLHD_OFF);
}

/**
 * @brief receive a packet in non-block method, it will return 0 when no data got
 * @param[in] <buff> buffer provide for data to receive
 * @return length, it will return 0 when no data got
 */
uint32_t rf_receive(uint8_t *buf)
{
	return PAN3028_recv_packet(buf);
}

/**
 * @brief receive a packet in non-block method, it will return 0 when no data got
 * @param[in] <buff> buffer provide for data to receive
			   <len> PLHD_LEN8 / PLHD_LEN16
 * @return result
 */
uint32_t rf_plhd_receive(uint8_t *buf,uint8_t len)
{
	return PAN3028_plhd_receive(buf,len);
}

/**
 * @brief rf enter rx continous mode to receive packet
 * @param[in] <none> 
 * @return result
 */
uint32_t rf_enter_continous_rx(void)
{
	if(PAN3028_set_mode(PAN3028_MODE_STB3) != OK)
	{
		return FAIL;
	}

	rf_port.set_rx();
	
	if(PAN3028_set_vco(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	}
	
	if(PAN3028_set_rx_mode(PAN3028_RX_CONTINOUS) != OK)
	{
		return FAIL;
	} 

	if(PAN3028_set_mode(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	} 
	return OK;
}

/**
 * @brief rf enter rx single timeout mode to receive packet
 * @param[in] <timeout> rx single timeout time(in ms)
 * @return result
 */
uint32_t rf_enter_single_timeout_rx(uint32_t timeout)
{
	if(PAN3028_set_mode(PAN3028_MODE_STB3) != OK)
	{
		return FAIL;
	}

	rf_port.set_rx();

	if(PAN3028_set_vco(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	}
	
	if(PAN3028_set_rx_mode(PAN3028_RX_SINGLE_TIMEOUT) != OK)
	{
		return FAIL;
	} 

	if(PAN3028_set_timeout(timeout) != OK)
	{
		return FAIL;
	}  

	if(PAN3028_set_mode(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	} 
	return OK;
}

/**
 * @brief rf enter rx single mode to receive packet
 * @param[in] <none> 
 * @return result
 */
uint32_t rf_enter_single_rx(void)
{
	if(PAN3028_set_mode(PAN3028_MODE_STB3) != OK)
	{
		return FAIL;
	}

	rf_port.set_rx();
	
	if(PAN3028_set_vco(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	}

	if(PAN3028_set_rx_mode(PAN3028_RX_SINGLE) != OK)
	{
		return FAIL;
	} 

	if(PAN3028_set_mode(PAN3028_MODE_RX) != OK)
	{
		return FAIL;
	} 
	return OK;
}

/**
 * @brief rf enter single tx mode and send packet
 * @param[in] <buf> buffer contain data to send
 * @param[in] <size> the length of data to send
 * @param[in] <tx_time> the packet tx time
 * @return result
 */
uint32_t rf_single_tx_data(uint8_t *buf, uint8_t size, uint32_t *tx_time)
{     
	if(PAN3028_set_mode(PAN3028_MODE_STB3) != OK)
	{
		return FAIL;
	}

	rf_port.set_tx();
	
	if(PAN3028_set_vco(PAN3028_MODE_TX) != OK)
	{
		return FAIL;
	}
	
	if(PAN3028_set_tx_mode(PAN3028_TX_SINGLE) != OK)
	{
		return FAIL;
	}  

	if(PAN3028_send_packet(buf, size) != OK)
	{
		return FAIL;
	}

	*tx_time = rf_get_tx_time();

	return OK;
}

/**
 * @brief rf enter continous tx mode to ready send packet
 * @param[in] <none> 
 * @return result
 */
uint32_t rf_enter_continous_tx(void)
{
	if(PAN3028_set_mode(PAN3028_MODE_STB3) != OK)
	{
		return FAIL;
	}

	rf_port.set_tx();
	
	if(PAN3028_set_vco(PAN3028_MODE_TX) != OK)
	{
		return FAIL;
	}
	
	if(PAN3028_set_tx_mode(PAN3028_TX_CONTINOUS) != OK)
	{
		return FAIL;
	}	

	return OK;
}
	
/**
 * @brief rf continous mode send packet
 * @param[in] <buf> buffer contain data to send
 * @param[in] <size> the length of data to send
 * @return result
 */
uint32_t rf_continous_tx_send_data(uint8_t *buf, uint8_t size)
{   
	if(PAN3028_send_packet(buf, size) != OK)
	{
		return FAIL;
	}

	return OK;
}

/**
 * @brief enable AGC function
 * @param[in] <state>  
 *			  AGC_OFF/AGC_ON
 * @return result
 */
uint32_t rf_set_agc(uint32_t state)
{
	if(PAN3028_agc_enable( state ) != OK)
	{
		return FAIL;
	}
	if(PAN3028_agc_config() != OK)
	{
		return FAIL;
	}
	return OK;
}

/**
 * @brief set rf para
 * @param[in] <para_type> set type, rf_para_type_t para_type
 * @param[in] <para_val> set value
 * @return result
 */
uint32_t rf_set_para(rf_para_type_t para_type, uint32_t para_val)
{
	PAN3028_set_mode(PAN3028_MODE_STB3);
	switch(para_type)
	{
		case RF_PARA_TYPE_FREQ:
			PAN3028_set_freq(para_val);  
			PAN3028_rst();
			break;
		case RF_PARA_TYPE_CR:
			PAN3028_set_code_rate(para_val);
			PAN3028_rst();
			break;
		case RF_PARA_TYPE_BW:
			PAN3028_set_bw(para_val);  
			PAN3028_rst();            
			break;
		case RF_PARA_TYPE_SF:
			PAN3028_set_sf(para_val);  
			PAN3028_rst();
			break;
		case RF_PARA_TYPE_TXPOWER:
			PAN3028_set_tx_power(para_val);
			PAN3028_rst(); 
			break;
		case RF_PARA_TYPE_CRC:
			PAN3028_set_crc(para_val);
			PAN3028_rst(); 
			break;
		default:
			break;    
	}
	return OK;
}

/**
 * @brief get rf para
 * @param[in] <para_type> get typ, rf_para_type_t para_type
 * @param[in] <para_val> get value
 * @return result
 */
uint32_t rf_get_para(rf_para_type_t para_type, uint32_t *para_val)
{
	PAN3028_set_mode(PAN3028_MODE_STB3);
	switch(para_type)
	{
		case RF_PARA_TYPE_FREQ:
			*para_val = PAN3028_read_freq();  
			break;
		case RF_PARA_TYPE_CR:
			*para_val = PAN3028_get_code_rate();
			break;
		case RF_PARA_TYPE_BW:
			*para_val = PAN3028_get_bw();
			break;
		case RF_PARA_TYPE_SF:
			*para_val = PAN3028_get_sf();          
			break;
		case RF_PARA_TYPE_TXPOWER:
			*para_val = PAN3028_get_tx_power();
			break;
		case RF_PARA_TYPE_CRC:
			*para_val = PAN3028_get_crc();          
			break;
		default:
			break;    
	}
	return OK;
}

/**
 * @brief set rf default para
 * @param[in] <none>
 * @return result
 */
void rf_set_default_para(void)
{
	PAN3028_set_mode(PAN3028_MODE_STB3);
	rf_set_para(RF_PARA_TYPE_FREQ, g_freq);
	rf_set_para(RF_PARA_TYPE_CR, DEFAULT_CR);
	rf_set_para(RF_PARA_TYPE_BW, g_band);
	rf_set_para(RF_PARA_TYPE_SF, g_factor);
	rf_set_para(RF_PARA_TYPE_TXPOWER, DEFAULT_PWR);
	rf_set_para(RF_PARA_TYPE_CRC, CRC_ON);

	rf_set_ldr(LDR_OFF);
}

/**
 * @brief set dcdc mode, The default configuration is DCDC_OFF, PAN3028 should set DCDC_OFF before enter sleep/deepsleep
 * @param[in] <dcdc_val> dcdc switch
 *		      DCDC_ON / DCDC_OFF
 * @return result
 */
uint32_t rf_set_dcdc_mode(uint32_t dcdc_val)
{
	return PAN3028_set_dcdc_mode(dcdc_val);
}

/**
 * @brief set LDR mode
 * @param[in] <mode> LDR switch
 *		      LDR_ON / LDR_OFF
 * @return result
 */
uint32_t rf_set_ldr(uint32_t mode)
{
	return PAN3028_set_ldr(mode);
}

/**
 * @brief set preamble by Spreading Factor,It is useful in all_sf_search mode
 * @param[in] <sf> Spreading Factor
 * @return result
 */
uint32_t rf_set_all_sf_preamble(uint32_t sf)
{
	return PAN3028_set_all_sf_preamble(sf);
}

/**
 * @brief open all sf auto-search mode
 * @param[in] <none> 
 * @return result
 */
uint32_t rf_set_all_sf_search(void)
{
	return PAN3028_set_all_sf_search( );
}

/**
 * @brief close all sf auto-search mode
 * @param[in] <none> 
 * @return result
 */
uint32_t rf_set_all_sf_search_off(void)
{
	return PAN3028_set_all_sf_search_off( );
}

/**
 * @brief set rf lna gain
 * @param[in] <mode>    
 *			  LNA_GAIN_LOW / LNA_GAIN_HIGH
 * @return result
 */
uint32_t rf_set_lna_gain(uint8_t mode)
{
	return PAN3028_set_lna_gain(mode);
}

/**
 * @brief RF IRQ handle process
 * @param[in] <none>
 * @return <none>
 */
void rf_irq_process(void)
{
    if(pan3028_irq_trigged_flag == true)
    {
		pan3028_irq_trigged_flag = false;
		
		uint8_t plhd_len;
		uint8_t irq = PAN3028_get_irq();

		if(irq & REG_IRQ_RX_PLHD_DONE)
		{
			plhd_len = PAN3028_get_plhd();
			rf_set_recv_flag(RADIO_FLAG_PLHDRXDONE);
			RxDoneParams.PlhdSize = PAN3028_plhd_receive(RxDoneParams.PlhdPayload,plhd_len);
			//PAN3028_rst();//stop it

		}
		if(irq & REG_IRQ_RX_DONE)
		{
			printf("In RX done\n");
#ifndef KERNEL
			RxDoneParams.Snr = PAN3028_get_snr();
			// printf("Before rssi\n");
			RxDoneParams.Rssi = PAN3028_get_rssi();
#endif
			rf_set_recv_flag(RADIO_FLAG_RXDONE);
			// printf("Before rcv package\n");
			RxDoneParams.Size = PAN3028_recv_packet(RxDoneParams.Payload);
			
		}
		if(irq & REG_IRQ_CRC_ERR)
		{
			rf_set_recv_flag(RADIO_FLAG_RXERR);
			PAN3028_clr_irq();

		}
		if(irq & REG_IRQ_RX_TIMEOUT)
		{
			PAN3028_rst();
			rf_set_recv_flag(RADIO_FLAG_RXTIMEOUT);
			PAN3028_clr_irq();

		}
		if(irq & REG_IRQ_TX_DONE)
		{
			rf_set_transmit_flag(RADIO_FLAG_TXDONE);
			PAN3028_clr_irq();

		}
	}
}

/**
 * @brief get one chirp time
 * @param[in] <bw>,<sf>
 * @return <time> us
 */
uint32_t get_chirp_time(uint32_t bw,uint32_t sf)
{
	switch(bw)
	{
		case 6:
			bw = 62500;
			break;
		case 7:
			bw = 125000;
			break;
		case 8:
			bw = 250000;
			break;
		case 9:
			bw = 500000;
			break;
		default:
			return FAIL;
	}
	return (1000000/bw)*(1<<sf);
}

/**
 * @brief check cad rx inactive
 * @param[in] <one_chirp_time>
 * @return <result> LEVEL_ACTIVE/LEVEL_INACTIVE
 */
/* uint32_t check_cad_rx_inactive(uint32_t one_chirp_time)
{
	uint8_t i = 0;

	rf_delay_us(one_chirp_time);
	rf_delay_us(360);
	for ( i = 0 ; i < 3 ; i++ )
    {
		rf_delay_us(one_chirp_time);
		if(CHECK_CAD() != 1)
		{
			if(rf_set_mode(PAN3028_MODE_STB3) != OK)
			{
				return FAIL;
			}
			return LEVEL_INACTIVE;
		}		
	}
    return LEVEL_ACTIVE;
} */

/**
 * @brief check cad tx inactive
 * @param[in] <none>
 * @return <result> OK/FAIL
 */
/* uint32_t check_cad_tx_inactive(void)
{
	uint32_t bw,sf;
	
	rf_get_para(RF_PARA_TYPE_BW, &bw);
	rf_get_para(RF_PARA_TYPE_SF, &sf);
	uint32_t one_chirp_time = get_chirp_time(bw,sf);//us

	if(rf_set_cad(CAD_DETECT_THRESHOLD_10) != OK)
	{
		return FAIL;
	}
	
	cad_tx_timeout_flag = MAC_EVT_TX_CAD_NONE;
	
	if(rf_enter_continous_rx() != OK)
	{
		return FAIL;
	}
	
	SET_TIMER_MS(one_chirp_time*3/1000 + 1);

	return OK;
}
 */


