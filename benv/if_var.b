#import "stdlib/core.b"
string root_source;

void init(){
	get_root_source(root_source); 
	SET_SOURCE(root_source); 
	SET_DEST("benv/if_var_program.b");
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

void main(){
	init();
	
	switch_command();
	
	#next_if:
	print("");
	if (NOT("end" == command)){
		if (is_if(command)){
			println(command);
		};
		switch_command();
		goto(#next_if);
	};

	finish();
};

main();
