int s_um (int a_g, int b_g){
	return (a_g + (b_g - 4));
};

bool bool_test(float a, float b, float c, float d){
	println("bool_test");
	return ((True)OR(False));
};

bool is_equal(){
	return ("lala=sdsd" == "lalsa=ldkf");
};
string b_str(){
	return "lalala";
};
int f_index(string a, string b){
	return 10;
};
void main(){
	string buf; 
	stack st;
	int t;
	push(5);
	pop(t);	

	st.push(5);
	st.pop(t); 
	
	buf = str(s_um(1, 2));
	buf = ((str(1) + str(2)) + b_str());
	buf = ((str(1) + "lalaal") + str(2));
	print(buf);
	print("\n");
	string s;
	s = "Hello world!";
	t = ((index(s, "na") + index(s,"ka")) + f_index("1", "2"));
	buf = str(t);
	print("\n");
	print(buf);
	print("\n");
};

main();
