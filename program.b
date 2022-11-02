#import "program2.b"
#import "program4.b"
#import "program5.b"


void main(){
	print("Hello world!!\n");
	int a;
	a = calc(1, 2);
	string buf;
	buf = str(a);
	print(buf);
	print("\n"); 
	a = sum2(3, 4);
	buf = str(a);
	print(buf);
	print("\n");
};

main();
