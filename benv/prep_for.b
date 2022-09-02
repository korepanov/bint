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

void main(){
	init();
	#next:
	switch_command();
	if (NOT("end" == command)){
		send_command(command);
		goto(#next);	
	};
	finish();
};

main();
