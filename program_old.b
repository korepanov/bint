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
	if (2 > 1){
		test_sum(sum(1, 2));
	}else{
		test_sum((sum(1, 2) + sum(3, 4)));
	};

	if (1 > 2){
		int a;
		print("YES");
		print("YES");
	}else if (3 > 4){
		print("YES");
		print("YES");
		float bb;	
	}else if (5 > 6){
		print("YES");
		bool ccc;
		print("YES");	
	};

	if (7 > 8){
		print("YES");
		stack dddd; 
		print("YES");	
	}else if (9 > 10){
		print("YES");
		int y;
		if (11 > 12){
			print("YES");
			string eeeee;
			bool ffffff;
			print("YES");
		};	
		float u;
	}else{
		println("OK");
	};

};

main(); 
