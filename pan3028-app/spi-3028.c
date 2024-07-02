#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/types.h>
#include <linux/spi/spidev.h>
#include <string.h>
#include <unistd.h>
#include <time.h>

#include "pan3028_port.h"
#include "spi-3028.h"
#include "radio.h"

static uint8_t bits = 8; 
static uint32_t speed = 100000; 
static uint16_t delay = 10;

int is_debug_enabled=0;

static const char *device = "/dev/spidev2.0";
static char customized_device[100] = {0,};

struct panchip3028_dev {
	int dummy;
};

struct panchip3028_dev panchip3028dev;

#define pr_info printf

void set_device(const char *device_name) {
	strncpy(customized_device, device_name, 90);
	customized_device[99] = '\0';
}
int g_freq = DEFAULT_FREQ;
int g_band = DEFAULT_BW;
int g_factor = DEFAULT_SF;

void set_freq(int freq)
{
	g_freq = freq;
}
void set_band(int band)
{
	g_band = band;
}
void set_factor(int factor)
{
	g_factor = factor;
}

static int fd = -1;

int open_dev()
{
	if(customized_device[0]==0) {
		strncpy(customized_device, device, 90);
	}
	// printf("Opening %s\n", customized_device);
	fd = open(customized_device, O_RDWR);

    if (fd < 0) { 
        printf("can't open device");
		// exit(-1);
	}
}

static int loop_started = 0;

int is_loop_started()
{
	return loop_started;
}

extern int zhg_usleep(int us);
int panchip3028_read_write_basic(struct panchip3028_dev *dev, void *tx_buf, void *rx_buf, int len)
{
	int ret; 

/* 	while (Reset == SPI_GetFlag(M4_SPI1, SpiFlagSendBufferEmpty))
	{
	}

	SPI_SendData8(M4_SPI1, tx_data);

	while (Reset == SPI_GetFlag(M4_SPI1, SpiFlagReceiveBufferFull))
	{
	}

	return SPI_ReceiveData8(M4_SPI1); */

	while(fd < 0)
	{
		open_dev();
	}
	
    struct spi_ioc_transfer tr = {
        .tx_buf = (unsigned long)tx_buf, 
        .rx_buf = (unsigned long)rx_buf, 
        .len = len, 
        .delay_usecs = delay, 
        .speed_hz = speed, 
        .bits_per_word = bits, 
    }; 
    //SPI_IOC_MESSAGE(1)
    ret = ioctl(fd, SPI_IOC_MESSAGE(1), &tr);
	// uint32_t mode;
	// ioctl(fd, SPI_IOC_RD_MODE, &mode);
	// printf("ioctl ret res: 0x%x, mode: 0x%x\n", ret, mode);

	return ret;
}

/**
 * @brief read one byte from register in current page
 * @param[in] <addr> register address to write
 * @return value read from register
 */
uint8_t PAN3028_read_reg(uint8_t addr)
{ 
	uint8_t temreg = 0x0;  
	uint8_t txdata[2] = {0, 0};
	uint8_t rxdata[2] = {0, 0};
	txdata[0] = (0x00 | (addr << 1));
	txdata[1] = 0x0;
	// printf("PAN3028_read_reg, 0x%x\n", addr);

	// rf_port.spi_cs_low();                               
	// rf_port.spi_readwrite(0x00 | (addr<<1));
	// temreg=rf_port.spi_readwrite(0x00);
	// rf_port.spi_cs_high();
	
	panchip3028_read_write_basic(&panchip3028dev, txdata, rxdata, 2);
	temreg = rxdata[1];
	// printf("PAN3028_read_reg, addr: 0x%x res_val: 0x%x\n", addr, temreg);
	return temreg;
}

uint32_t PAN3028_write_reg(uint8_t addr,uint8_t value)
{ 
	uint16_t tmpreg = 0;  
	uint16_t addr_w = (0x01 | (addr << 1));
	uint8_t txdata[2] = {0, 0};
	uint8_t rxdata[2] = {0, 0};
	txdata[0] = addr_w;
	txdata[1] = value;
	// printf("PAN3028_write_reg 0x%x, val: 0x%x\n", addr, value);

	// rf_port.spi_cs_low();	  
	// rf_port.spi_readwrite(addr_w);
	// rf_port.spi_readwrite(value);
	// rf_port.spi_cs_high();
	panchip3028_read_write_basic(&panchip3028dev, txdata, rxdata, 2);

	// printf("PAN3028_write_reg write complete\n");
	tmpreg = PAN3028_read_reg(addr);
	if (is_debug_enabled){
		printf("PAN3028_write_reg read from 0x%x again, target: 0x%x, val: 0x%x\n", addr, value, tmpreg);
	}

	if(tmpreg == value)
	{
		return OK;
	}
	else
	{
		return FAIL;
	}	
}


static int panchip3028_readwritebyte(struct panchip3028_dev *dev, uint8_t tx_data_byte)
{
	int ret = -1;
	unsigned char txdata[1] = {0, };
	unsigned char rxdata[1] = {0, };
	int res = 

	txdata[0] = tx_data_byte;
	
	res = panchip3028_read_write_basic(dev, txdata, rxdata, 1);
	
	if (res < 0) {
		return res;
	}

	ret = rxdata[0];
	
	return ret;
}

uint8_t spi_readwritebyte(uint8_t tx_data){
    int res = panchip3028_readwritebyte(&panchip3028dev, tx_data);
    if (res < 0) {
        pr_info("Read write error in spi for data: 0x%x", tx_data);
        return 0;
    }
    return res;
}

void dump_data(unsigned char *buf, int len)
{
	int i;
	for(i = 0; i < len; i++)
	{
		pr_info("0x%02x ", buf[i]);
	}
}

int spi_readwrite_reg_with_buf(uint8_t reg_byte, void* tx_data, void* rx_data, int len)
{
	uint8_t tx[300];
	uint8_t rx[300];
	int res = -1;
	if(len>299){
		pr_info("spi_readwrite_reg_with_buf only support len < 300, 299 byte data + 1 reg byte\n");
		return -1;
	}
	tx[0] = reg_byte;
	if(tx_data == NULL){
		// printf("tx buf: 0x%x\n", tx);
		memset(tx+1, 0, len);
	}
	else{
		memcpy(tx+1, tx_data, len);
	}
	pr_info("Writing buffer, len: %d\n", len);
	dump_data(tx, len+1);
	pr_info("\n");
	res = panchip3028_read_write_basic(&panchip3028dev, tx, rx, len+1);
	if(rx_data != NULL) {
		// Ignore the data get when sending reg_byte
		memcpy(rx_data, rx+1, len);
	}
	return res;
}


int zhg_usleep(int us){
#ifdef KERNEL
	usleep_range(us, us+1);
#else
	usleep(us);
#endif
}

#define DDL_Printf printf
#define SysTick_Delay zhg_usleep

extern struct RxDoneMsg RxDoneParams;
#define MAX_SENDING_LEN 512
char tx_test_buf[MAX_SENDING_LEN] = {0,};
int tx_len = 0;
uint32_t tx_time=0;
volatile int is_tx_buf_loaded = 0;

char rx_buf[MAX_SENDING_LEN] = {0,};
volatile int received_len = 0;

int send(unsigned char* buf, int len)
{
	if(len>=MAX_SENDING_LEN){
		pr_info("Can only send %d bytes\n", MAX_SENDING_LEN);
		return -1;
	}
	if(is_tx_buf_loaded){
		pr_info("is still sending data, can not send again\n");
		return -1;
	}
	memcpy(tx_test_buf, buf, len);
	tx_len = len;
	is_tx_buf_loaded = 1;
	return 0;
}

volatile int is_rx_buf_empty = 1;

int receive(unsigned char* buf, int buffer_len)
{
	// pr_info("RX: buf: 0x%x len: %d\n", buf, buffer_len);
	int res = 0;

	if(is_rx_buf_empty){
		// pr_info("RX buf is empty\n");
		return res;
	}
	if(buffer_len < received_len){
		pr_info("RX buf is too small\n");
		return -1;
	}
	pr_info("memcpy buf: %x\n, recvd len: %d\n", buf, received_len);
	memcpy(buf, rx_buf, received_len);
	is_rx_buf_empty = 1;
	res = received_len;
	received_len = 0;
	return res;
}

int is_in_rx_mode = 0;
int is_tx_ongoing = 0;
time_t tx_timestamp;

int init_tx_or_rx()
{
	if(is_tx_buf_loaded)
	{
		if(rf_single_tx_data(tx_test_buf, tx_len, &tx_time) != OK)
		{
			DDL_Printf("tx fail \r\n");
		}
		else
		{
			is_tx_buf_loaded = 0;
		}
		is_in_rx_mode = 0;
		is_tx_ongoing = 1;
		time (&tx_timestamp);
		DDL_Printf("Send initated, %d\n", tx_timestamp);
	}
	else
	{//TODO: Remove this RX as after TX, we will call it
		time_t t;
		time (&t);
		if (is_tx_ongoing && (t-tx_timestamp)>5) {
			is_tx_ongoing = 0;
			DDL_Printf("Tx timeout, start: %d, %d\n", tx_timestamp, t);
		}
        if(!is_in_rx_mode && !is_tx_ongoing){
			pr_info("Switch to continous RX mode\n");
			rf_enter_continous_rx();
			is_in_rx_mode = 1;
		}
		
	}
}

double RSSI;
double SNR;

void fill_rx_buf()
{
	if(!is_rx_buf_empty) {
		pr_info("Buffer is not empty, but new data received\n");
	}
	else {
		if (RxDoneParams.Size >= 255){
			pr_info("RxDoneParams size overflow\n");
		}
		is_rx_buf_empty = 0;
		received_len = RxDoneParams.Size;
		memcpy(rx_buf, RxDoneParams.Payload, RxDoneParams.Size);
		RSSI = (double)RxDoneParams.Rssi;
		SNR = (double)RxDoneParams.Snr;
	}
}

static uint32_t txcnt = 0;
static uint32_t rxcnt = 0;

void event_handler()
{
	uint8_t i;
	is_debug_enabled = 0;
	rf_irq_handler();
	rf_irq_process();
	//rxdone flag
	// DDL_Printf("Send receive loop, before rf_get_recv_flag\n");
	if(rf_get_recv_flag() == RADIO_FLAG_RXDONE)
	{
		// BSP_LED_Toggle();
		rf_set_recv_flag(RADIO_FLAG_IDLE); 
		DDL_Printf("Rx : SNR: %f ,RSSI: %f \r\n", RxDoneParams.Snr, RxDoneParams.Rssi);
		for(i = 0; i < RxDoneParams.Size; i++)
		{
			DDL_Printf("0x%02x ", RxDoneParams.Payload[i]);
		}
		fill_rx_buf();
		DDL_Printf("\r\n");
		rxcnt ++;
		DDL_Printf("###Rx cnt %d##\r\n", rxcnt);
		//rxdone,set sleep and wakeup
		rf_sleep();
		rf_sleep_wakeup();

		SysTick_Delay(3000);
		//tx
		// send_and_update_flag();
		is_in_rx_mode = 0;

	}
	// DDL_Printf("Send receive loop, before rf_get_recv_flag timeout and error\n");
	//rxtimeout or rxerr flag
	if((rf_get_recv_flag() == RADIO_FLAG_RXTIMEOUT) || (rf_get_recv_flag() == RADIO_FLAG_RXERR))
	{
		rf_set_recv_flag(RADIO_FLAG_IDLE); 
		DDL_Printf("Rxerr\r\n");
		//rxtimeout or rxerr, set sleep and wakeup
		rf_sleep();
		rf_sleep_wakeup();
		//tx again
		SysTick_Delay(10000);
		// send_and_update_flag();
		is_in_rx_mode = 0;
	}
	// DDL_Printf("Send receive loop, before rf_get_transmit_flag\n");
	if(rf_get_transmit_flag() == RADIO_FLAG_TXDONE)
	{
		rf_set_transmit_flag(RADIO_FLAG_IDLE);  
		txcnt ++;
		DDL_Printf("Tx cnt %d\r\n", txcnt );			
		//txdone, set sleep and wakeup
		rf_sleep();
		rf_sleep_wakeup();
		is_tx_ongoing = 0;
		// rf_enter_continous_rx();
	}
}

void toggle_debug()
{
	is_debug_enabled = !is_debug_enabled;
}