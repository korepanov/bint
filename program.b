int s_um (int a_g, int b_g){
	return (a_g + (b_g - 4));
};

bool bool_test(float a, float b, float c, float d){
	print("bool_test");
	print("\n");
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
	stack poses;
	int beg;
	int t;
	bool b;
	string command;
	string root_source;
	string root_dest;
	
	[goto(#mark1), ("yes" == buf[1:5]), print("")];
	send_command(buf[1:5]);
	RESET_SOURCE();
	UNDEFINE(buf);
	#mark2:
	goto(#mark1);
	#mark1:
	goto(#mark2);
	get_root_source(root_source);
	get_root_dest(root_dest);
	SET_SOURCE("program_old.b");
	SET_DEST("p.b");
	next_command(command);
	send_command(command);
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST("temp.b");
	SEND_DEST("p.basm");
	REROUTE();

	input(buf);
	print(buf);
	print("\n");
	
	st = reg_find("a","b");
	st.pop(buf);
	print(buf);
	poses.pop(beg);
	buf = str(beg);
	print(buf);
	print("\n");

	b = is_letter("f");
	buf = str(b);
	print(buf);
	print("\n");
	b = ((is_digit("a"))AND(is_letter("b")));
	buf = str(b);
	print(buf);
	print("\n");
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
