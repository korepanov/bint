#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_if_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_if(string command){
	string op;
	stack s;
	string buf;
	int ibuf;
	int start_pos;

	op = "if(";	
	s = ops(command, op);
	s.pop(buf);
	[goto(#end_true), ("end" == buf), print("")];
	ibuf = int(buf);
	[goto(#start_true), (0 == ibuf), print("")];
	start_pos = (ibuf - 4);
	[goto(#end_true), (start_pos < 0), print("")];
	buf = command[start_pos:ibuf];
	[print(""), ("else" == buf), goto(#end_true)];
	#start_true:
	return True;
	#end_true:
	return False;
	print("");
};

bool is_else(string command){
	string op;
	stack s;
	string buf;
	int ibuf;
	int start_pos;

	op = "}else{";
	s = ops(command, op); 
	s.pop(buf);
	[goto(#end_else_true), ("end" == buf), print("")];
	ibuf = int(buf);
	[print(""), (0 == ibuf), goto(#end_else_true)];
	return True;
	#end_else_true:
	return False;
};

string add_nop(string command){
	string op;
	string buf;
	string command_buf; 
	stack s; 
	int pos;
	string new_command;
	int command_len;
	int num;
	
	op = "{";
	num = 0;

	#add_nop_s:
	s = ops(command, op);
	s.pop(buf);
	[print(""), ("end" == buf), goto(#add_nop_e)];
	
	pos = int(buf);
	pos = (pos + 1);
	command_buf = command[0:pos];
	new_command = (command_buf + "print(\"\")");
	command_len = len(command);
	command_buf = command[pos:command_len];
	new_command = (new_command + command_len);
	command = new_command; 

	#add_nop_e:
	return "plug";
};

void main(){
	string command;

	init();
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	[print(""), ((is_if(command)) OR (is_else(command))), goto(#not_cond)];
	println(add_nop(command));
	#not_cond:
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
