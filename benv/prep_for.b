#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_for_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_for(string command){
	stack s;
	string buf;

	buf = "for(";
	s = ops(command, buf);
	s.pop(buf);
	
	return (NOT("end" == buf));
};

void main(){
	init();
	#next:
	switch_command();
	if (NOT("end" == command)){
		if (is_for(command)){
			println(command);		
		};
		send_command(command);
		goto(#next);	
	};
	finish();
};

main();
