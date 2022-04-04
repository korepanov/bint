string command;
bool bool_res;

int init(){
	SET_SOURCE("benv/prep_func_program.b");
	SET_DEST("benv/long_function_program.b");
	next_command(command);

	return 0;
};

int finish(){
	DEL_DEST("benv/prep_func_program.b");
	UNSET_SOURCE();
	UNSET_DEST();

	return 0;
};

bool is_func(string command){
	int number; 
	number = index(command, "{");
	[goto(#no), (-1 == number), print("")];
	return True;
	#no:
	return False;
};

int res; 
res = init();

[print(""), (0 == res), print("INIT ERROR\n")];

#start:
[goto(#end), ("end" == command), print("")];
send_command(command);
bool_res = is_func(command);
[print(command), (bool_res), print("")];
next_command(command);
goto(#start);
#end:
res = finish();
[print(""), (0 == res), print("FINISH ERROR\n")];
