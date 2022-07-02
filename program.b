int sum (int a, int b){
	return (a + b);
};


void main(){
	string buf; 
	buf = str(sum(1, 2));
	buf = (str(1) + str(2));
	buf = ((str(1) + "lalaal") + str(2));
	print(buf);
	print("\n");
};

main();
