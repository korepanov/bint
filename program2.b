#import "program3.b"
#import "program4.b" 

int calc(int a, int b){
	int c;
	int d;
	c = sum(a, b);
	d = sum2(a, b);
	c = sum(c, d);
	return c; 
};
