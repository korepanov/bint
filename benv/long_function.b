string command;

int init(int a){
	SET_SOURCE("benv/prep_func_program.b");
	SET_DEST("benv/long_function_program.b");
	next_command(command);

	return 0;
};

int res; 
res = init(5);

[print(""), (0 == res), print("INIT ERROR\n")];

#start:
[goto(#end), ("end" == command), print("")];
send_command(command);
next_command(command);
goto(#start);
#end:
DEL_DEST("benv/prep_func_program.b");
UNSET_SOURCE();
UNSET_DEST(); 
