#define pr_info printf
extern int zhg_usleep(int us);
extern void set_device(const char *device_name);
extern int g_freq;
extern int g_band;
extern int g_factor;
void set_freq(int freq);
void set_band(int band);
void set_factor(int factor);
int receive(unsigned char* buf, int buffer_len);
int send(unsigned char* buf, int len);
void toggle_debug();
int is_loop_started();
void event_handler();
int init_tx_or_rx();
extern double RSSI;
extern double SNR;
