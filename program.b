int s_um (int a_g, int b_g){
	return (a_g + b_g);
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
