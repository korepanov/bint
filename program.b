float sum(float a, float b){
	float res; 
	res = (a + b); 
	return res;
};

int res_hello(){
	print("Hello, world!\n");
	return 0; 
};

void hello(string user_name){
	string buf; 
	buf = (("Hello, " + user_name) + "!\n");
	print(buf);
};

void hello_world(){
	print("Hello, world!\n");
};



string buf;
print("Input your name\n"); 
input(buf);
hello(buf);
float c; 
c = sum(5, 6);
string buf; 
buf = str(c); 
print(buf);
print("\n");
int res; 
res = res_hello(); 
[print("OK\n"), (0 == res), print("ERROR\n")];
hello_world(); 
