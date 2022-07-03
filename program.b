int s_um (int a_g, int b_g){
	return (a_g + b_g);
};

bool bool_test(float a, float b, float c, float d){
	println("bool_test");
	return (((NOT(1==1))AND(NOT(1>2)))XOR(arg));
};

bool is_equal(){
	return ("lala=sdsd" == "lalsa=ldkf");
};
string b_str(){
	return "lalala";
};
void main(){
	string buf; 
	buf = str(s_um(1, 2));
	buf = ((str(1) + str(2)) + b_str());
	buf = ((str(1) + "lalaal") + str(2));
	print(buf);
	print("\n");
};

main();
