#import "stdlib/core.b"
bool first_file; 

void init(){
	string root_source;
	first_file = True;
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

int func_call(string fname, string command){
	string reg;
	stack s; 
	stack this_s;
	string buf;
	int pos;
	
	reg = (("(?:" + fname) + "\(.*\))");
	s = reg_find(reg, command);
	s.pop(this_s);
	this_s.pop(buf);
	
	if ("end" == buf){
		pos = -1;
	}else{
		pos = int(buf);
	};
	
	
	return pos;
};

void switch_files(){
	finish();
	[print(""), (first_file), goto(#first_end)];
	SET_SOURCE("benv/recurs_program.b");
	SET_DEST("benv/recurs_program2.b");
	first_file = False;
	goto(#switch_files_e);
	#first_end:
	SET_SOURCE("benv/recurs_program2.b");
	SET_DEST("benv/recurs_program.b");
	first_file = True; 
	#switch_files_e:
	print("");
};

void clear_files(){
	[goto(#first_file), (first_file), print("")];
	switch_files();
	switch_command();
	#clear_files_s:
	[goto(#clear_files_e), ("end" == command), print("")];
	send_command(command);
	switch_command();
	goto(#clear_files_s);

	#first_file:
	print("");
	#clear_files_e:
	finish();
	DEL_DEST("benv/recurs_program2.b");
};


void main(){
	init();
	string name;
	stack accepted_args;
	stack raccepted_args;
	string t;
	stack arg;
	string arg_name;
	string arg_type;
	int counter;
	string buf;
	int old_COMMAND_COUNTER;
	int pos;
	int call_counter;
	int old_call_counter;
	string new_command;
	bool was_recurs;
	bool first_recurs_call;
	
	call_counter = 0;
	was_recurs = False;
	first_recurs_call = True;
	
	switch_command();
	
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
			
			raccepted_args = args_to_accept(command);
 			
			counter = block_end();
			old_COMMAND_COUNTER = COMMAND_COUNTER;
			
			#is_recurs:
			print("");
			if (COMMAND_COUNTER < counter){
				switch_command();
				pos = func_call(name, command);
				if (NOT(-1 == pos)){
					was_recurs = True;
					old_call_counter = call_counter;
					call_counter = COMMAND_COUNTER;
					SET_COMMAND_COUNTER(old_call_counter);
					
					switch_command();
					#block_begin:
					print("");
					if (COMMAND_COUNTER < old_COMMAND_COUNTER){
						send_command(command);
						switch_command();
						goto(#block_begin);
					};
					
					if (first_recurs_call){
						first_recurs_call = False;
						new_command = ("stack" + ("$" + (name + "_stack")));
						send_command(new_command);
						
						new_command = (("#" + name) + ":");
						new_command = (new_command + "print(\"\")");
						send_command(new_command); 
						switch_command();
					};
					
					#before_recurs:
					print("");
					if (COMMAND_COUNTER < call_counter){
						send_command(command);
						switch_command();
						goto(#before_recurs);
					};
					
					accepted_args = raccepted_args;
					accepted_args.pop(arg);

					arg.pop(arg_type);
					arg.pop(arg_name);

					#args:
					print("");
					if (NOT("end" == arg_name)){
						new_command = (("push(" + arg_name) + ")");
						send_command(new_command);
						
						accepted_args.pop(arg);
						arg.pop(arg_type);
						arg.pop(arg_name);
				
						goto(#args);
					};
					
					
				};
				goto(#is_recurs);
			};
			
			if (was_recurs){
				was_recurs = False; 
				first_recurs_call = True;
				
				SET_COMMAND_COUNTER(call_counter);
				
				switch_command();
				
				#rest:	
				print("");
				if (COMMAND_COUNTER < counter){
					send_command(command);
					switch_command();
					goto(#rest);
					
				}; 
				send_command(command);
				
				switch_command();
				
				#external_call:
				print("");
				if(NOT("end" == command)){
					pos = func_call(name, command);
					if (NOT(-1 == pos)){
						println(command);
					};
					send_command(command);
					switch_command();
					goto(#external_call);
				};
				switch_files();
				call_counter = 0;
				SET_COMMAND_COUNTER(0);
			};
			
		};
		switch_command();
		goto(#main_s);
	};
	
	switch_command();
	
	#final_rest:
	if (NOT("end" == command)){
		send_command(command);
		switch_command();
		goto(#final_rest);
	};
	
	new_command = "#_exit:print(\"\")";
	send_command(new_command);
	clear_files(); 
	
};

main();


