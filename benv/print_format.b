#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/print_format_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void replace_print(string command){
	stack s;
	string buf;
	int nbuf;
	int arg_begin;
	int arg_end;
	int num;
	string snum;
	int command_len;
	int i;
	string op;

	op = "print(";
	i = 0;
	command_len = len(command);
	num = 0;
	snum = str(num);
	s = ops(command, op);

	s.pop(buf);
	#replace_s:
	[goto(#replace_e), ("end" == buf), print("")];
	nbuf = int(buf);
	nbuf = (nbuf + 5);
	buf = str(nbuf);
	arg_begin = (nbuf + 1);
	arg_end = func_end(command, nbuf);
	buf = str(nbuf);
	buf = ("string print_arg" + snum);
	send_command(buf);
	buf = command[arg_begin:arg_end];
	buf = ((("print_arg" + snum) + "=") + buf);
	send_command(buf);
	buf = command[arg_end:command_len]; 
	command = command[0:arg_begin];
	command = ((command + "print_arg") + snum);
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
	buf = (("UNDEFINE(print_arg" + snum) + ")");
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
