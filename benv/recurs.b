void init(){
	string root_source;
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/recurs_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_func_definition(string command){
	stack s;
	string buf;

	string r1;
	string r2;
	string r;

	r1 = "(?:(int|float|bool|string|stack|void)[[:alnum:]|_]";
	r2 = "*?\((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})(int|float|bool|string|stack)[[:alnum:]|_]+\){";
	r = (r1 + r2);

	s = reg_find(r, command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return True;
	};	

	s = reg_find("(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\(\){)", command);
	
	s.pop(buf);
	return (NOT("end" == buf));
};

void main(){
	init();
	string command;
	next_command(command);
	
	#main_s:
	print("");
	if (NOT("end" == command)){
		if (is_func_definition(command)){
			print(command);
			print("\n");		
		};
		send_command(command);
		next_command(command);
		goto(#main_s);
	};
	finish();
};

main();


