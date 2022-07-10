#import "stdlib/core.b"

string root_dest;

void init(){
	SET_SOURCE("benv/func_program.b");
	get_root_dest(root_dest);
	SET_DEST("benv/program.basm");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	DEL_DEST("benv/func_program.b");
	SEND_DEST(root_dest);
};

void replace_print(string command){
	stack s;
	string buf;
	string buf2;
	int nbuf;
	int arg_begin;
	int arg_end;
	int num;
	string snum;
	int command_len;
	int i;
	string op;
	string first_symbol;

	first_symbol = command[0];
	[print(""), ("[" == first_symbol), goto(#send_this_end)];
	#send_this:
	send_command(command);
	goto(#un_loop_e);
	#send_this_end:
	op = "print(";
	i = 0;
	command_len = len(command);
	num = 0;
	snum = str(num);
	print(command);
	print("\n");
	print(op);
	print("\n");
	s = ops(command, op);
	
	s.pop(buf);
	#replace_s:
	[goto(#replace_e), ("end" == buf), print("")];
	nbuf = int(buf);
	nbuf = (nbuf + 5);
	buf = str(nbuf);
	arg_begin = (nbuf + 1);
	arg_end = func_end(command, nbuf);
	buf2 = command[arg_begin:arg_end];
	[goto(#send_this), ("\"\"" == buf2), print("")];
	buf = str(nbuf);
	buf = ("string $print_arg" + snum);
	send_command(buf);
	buf = command[arg_begin:arg_end];
	buf = ((("$print_arg" + snum) + "=") + buf);
	send_command(buf);
	buf = command[arg_end:command_len]; 
	command = command[0:arg_begin];
	command = ((command + "$print_arg") + snum);
	command = (command + buf);
	command_len = len(command);
	s = ops(command, op);
	num = (num + 1);
	snum = str(num);
	
	#loop:
	[print(""), (i < num), goto(#loop_e)];
	s.pop(buf);
	i = (i + 1);	
	goto(#loop);
	#loop_e:
	i = 0;
	s.pop(buf);	
	goto(#replace_s);
	#replace_e:	
	send_command(command);
	i = 0;
	#un_loop:
	[print(""), (i < num), goto(#un_loop_e)];
	snum = str(i);
	buf = (("UNDEFINE($print_arg" + snum) + ")");
	send_command(buf);
	i = (i + 1);
	goto(#un_loop);
	#un_loop_e:
	print("");
};

void main(){
	string command;
	init();
	
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	replace_print(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
