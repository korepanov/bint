int sum(int a, int b){
	return (a + b);
};

void test_sum(string t){
	print(t);
	print("\n");
};

void main(){
	test_sum(sum(1, 2));

};

main(); 
