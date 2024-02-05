#import "stdlib/core.b"

string root_source;
int num;
bool first_file;
string condition;

void init(){
	num = 0;
	first_file = True; 
	root_source = "benv/while_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/for_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_for(string command){
	stack s;
	string buf;
	int pos;

	buf = "for(";
	s = ops(command, buf);
	s.pop(buf);

	if (NOT("end" == buf)){
		pos = int(buf);
		if (NOT(0 == pos)){
			println("for: ERROR");
			exit(1);	
		};	
	};

	return (NOT("end" == buf));
};

bool is_var_def(string command){
	stack s;
	string buf;

	s = reg_find("^(?:(int|float|bool|stack|string)[^\(]*)", command);
	s.pop(buf);
	
	return (NOT("end" == buf)); 
};

string Type(string command){
	stack s;
	string buf;

	s = reg_find("(?:(^int))", command);
	s.pop(buf);
	
	[goto(#int_end), ("end" == buf), print("")];
	return "int";	
	#int_end:

	s = reg_find("(?:(^float))", command);
	s.pop(buf);
	[goto(#float_end), ("end" == buf), print("")];
	return "float";	
	#float_end:

	s = reg_find("(?:(^bool))", command);
	s.pop(buf);
	[goto(#bool_end), ("end" == buf), print("")];
	return "bool";	
	#bool_end:

	s = reg_find("(?:(^stack))", command);
	s.pop(buf);
	[goto(#stack_end), ("end" == buf), print("")];
	return "stack";	
	#stack_end:

	s = reg_find("(?:(^string))", command);
	s.pop(buf);
	[goto(#string_end), ("end" == buf), print("")];
	return "string";	
	#string_end:

	s = reg_find("(?:(^void))", command);
	s.pop(buf);
	[goto(#void_end), ("end" == buf), print("")];
	return "void";	
	#void_end:
	
	print(command);
	print("\n");
	print("Type: ERROR\n");
	exit(1);
};

void switch_files(){
	finish();
	[print(""), (first_file), goto(#first_end)];
	SET_SOURCE("benv/for_program.b");
	SET_DEST("benv/for_program2.b");
	first_file = False;
	goto(#switch_files_e);
	#first_end:
	SET_SOURCE("benv/for_program2.b");
	SET_DEST("benv/for_program.b");
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
	DEL_DEST("benv/for_program2.b");
	DEL_DEST(root_source);
};

void main(){
	init();
	int counter;
	int internal_counter;
	int command_len;
	int command_len2;
	int buf_len;
	int old_num;
	string snum;
	string buf;
	string inc;
	string old_command;
	string internal_old_command; 
	int pos;
	bool was_for;
	bool was_internal_for;
	string arg_type;
	string arg_name; 
	
	
	was_for = False;
	was_internal_for = False;
	
	#next:
	switch_command();
	if (NOT("end" == command)){

		if (is_for(command)){
			was_for = True;
			was_internal_for = False;
			command_len = len(command);
			command = command[4:command_len];
			send_command(command);
			old_command = command;
			switch_command();
			send_command(command);
			switch_command();
			snum = str(num);
			buf = ("bool $for" + snum);
			send_command(buf);
			buf = (("#_for" + snum) + ":print(\"\")");
			send_command(buf);  
			buf = (("if(" + command) + "){print(\"\")");
			condition = buf;
			send_command(buf);
			switch_command();
			counter = block_end();
			pos = index(command, "{");
			if (-1 == pos){
				println("for: ERROR");
				exit(1);			
			};
			pos = (pos - 1);
			inc = command[0:pos];
			#next_internal:
			switch_command(); 
			if (COMMAND_COUNTER < counter){
				if (is_for(command)){
					was_internal_for = True;
					send_command(command); 
					switch_command();
					send_command(command);
					switch_command();
					send_command(command);
					switch_command();
					send_command(command);	
					internal_counter = block_end();	
					switch_command(); 			
				};
				
				if (was_internal_for){
					if (COMMAND_COUNTER > internal_counter){
						was_internal_for = False; 					
					};				
				};

				if (("break" == command) AND (NOT(was_internal_for))){
					command = (("$for" + snum) + "=True"); 
					send_command(command);
					command = (("goto(#_undef_for" + snum) + ")");				
				};
				if (("continue" == command) AND (NOT(was_internal_for))){
					send_command(inc);
					command = (("goto(#_undef_for" + snum) + ")");					
				};
				send_command(command);
				goto(#next_internal);			
			};
			send_command(inc);
			buf = (("#_undef_for" + snum) + ":print(\"\")");
			send_command(buf);
			send_command(command);
			buf = (("if($for" + snum) + "){print(\"\")");
			send_command(buf);
			buf = ((("goto(#_for") + snum) + "_end)");
			send_command(buf);
			buf = "}";
			send_command(buf);			
			send_command(condition);
			buf = "CLEAR()";
			send_command(buf); 
			buf = (("goto(#_for" + snum) + ")");
			send_command(buf);
			buf = "}";
			send_command(buf); 
			command = (("#_for" + snum) + "_end:print(\"\")");
			send_command(command); 
			command = (("UNDEFINE($for" + snum) + ")");
			send_command(command);
			
			if(is_var_def(old_command)){
				arg_type = Type(old_command);
				int type_len;
				type_len = len(arg_type);
				command_len = len(old_command);
				arg_name = old_command[type_len:command_len];
				old_command = (("UNDEFINE(" + arg_name) + ")");
				send_command(old_command); 				
			};

			old_num = num;
			num = (num + 1);
		}else{
			send_command(command);
		};
		goto(#next);	
	};

	if (was_for){
		was_for = False;
		switch_files();
		COMMAND_COUNTER = 0;
		goto(#next); 
	};

	clear_files(); 
};

main();
