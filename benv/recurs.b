#import "stdlib/core.b"

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
	string r3;
	string r;

	r1 = "(?:(int|float|bool|string|stack|void)[[:alnum:]|_]";
	r2 = "*?\((?:((int|float|bool|string|stack))[[:alnum:]|_]+\,){0,})";
	r3 = "(int|float|bool|string|stack)[[:alnum:]|_]+\){";
	
	r = ((r1 + r2) + r3);

	s = reg_find(r, command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return True;
	};	

	s = reg_find("(?:(int|float|bool|string|stack|void)[[:alnum:]|_]*?\(\){)", command);
	
	s.pop(buf);
	return (NOT("end" == buf));
};

string T(string command){
	stack s;
	string buf;

	s = reg_find("(?:(^int))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "int";	
	};

	s = reg_find("(?:(^float))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "float";	
	};

	s = reg_find("(?:(^bool))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "bool";	
	};

	s = reg_find("(?:(^stack))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "stack";	
	};

	s = reg_find("(?:(^string))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "string";	
	};

	s = reg_find("(?:(^void))", command);
	s.pop(buf);
	if (NOT("end" == buf)){
		return "void";	
	};
	
	print(command);
	print("\n");
	print("T: ERROR\n");

	return "ERROR";	
};

string func_name(string command){
	string t;
	int pos;
	int epos;
	int command_len;

	t = T(command);
	if ("ERROR" == command){
		return "ERROR";	
	};
	
	pos = len(t);
	command_len = len(command);
	command = command[pos:command_len];
	
	epos = index(command, "(");

	if (-1 == epos){
		return "ERROR";	
	};

	command = command[0:epos];

	return command;

};

stack args_to_accept(string command){
	int pos1;
	int pos2;
	string s_to_parse;
	string t;
	stack res;
	stack null;
	stack el;
	int L;
	string arg_name;
	string arg_type;
	
	pos1 = index(command, "(");
	pos1 = (pos1 + 1);
	pos2 = index(command, ")");

	if (pos1 == pos2){
		return res; 	
	};

	if ((-1 == pos1) OR (-1 == pos2)){
		print(command);
		print("\n");
		print("recurs: args_to_accept: ERROR\n");	
	};

	s_to_parse = command[pos1:pos2];
	
	#parse_s:	
	el = null;
	t = T(s_to_parse);
	pos1 = len(t);
	L = len(s_to_parse);

	s_to_parse = s_to_parse[pos1:L];
	
	pos1 = index(s_to_parse, ",");
	if (-1 == pos1){
		pos1 = len(s_to_parse);	
		arg_name = s_to_parse[0:pos1];
		el.push(arg_name);
		el.push(t);
		res.push(el);

		stack res_el;
		stack new_res;
		string buf;

		#rev:
		res.pop(res_el);
		res_el.pop(buf);

		if ("end" == buf){
			return new_res;		
		};

		res_el.push(buf);
		new_res.push(res_el);		

		goto(#rev); 
	};
	
	arg_name = s_to_parse[0:pos1];
	pos1 = (pos1 + 1);
	L = len(s_to_parse);
	s_to_parse = s_to_parse[pos1:L];
	el.push(arg_name);
	el.push(t);
	res.push(el);

	goto(#parse_s);
};

void main(){
	init();
	string command;
	string name;
	stack accepted_args;
	string t;
	stack arg;
	string arg_name;
	string arg_type;

	next_command(command);
	
	#main_s:
	print("");
	if (NOT("end" == command)){
		if (is_func_definition(command)){
			t = T(command);
			if ("ERROR" == t){
				print("recurs ERROR\n");			
			};
			name = func_name(command);
			if ("ERROR" == name){
				print("recurs ERROR\n");			
			};
			accepted_args = args_to_accept(command);
 			
			accepted_args.pop(arg);

			arg.pop(arg_type);
			print(name);
			print("\n");

			if ("end" == arg_type){
				print("no args\n");			
			};
			
			arg.pop(arg_name);

			#args:
			if (NOT("end" == arg_name)){
				print(arg_type);
				print(" ");
				print(arg_name);
				print("\n");
				
				accepted_args.pop(arg);
				arg.pop(arg_type);
				arg.pop(arg_name);
				
				goto(#args);
			};

		};
		send_command(command);
		next_command(command);
		goto(#main_s);
	};
	finish();
};

main();


