#import "stdlib/core.b"

int sum(int a, int b){
	return (a + b);
};

void test_sum(int t){
	string buf;
	buf = str(t);
	print(buf);
	print("\n");
};

void main(){
	println("Hello world!");
	test_sum(sum(1, 2));
	test_sum((sum(1, 2) + sum(3, 4)));

};

main(); 
