$filestdlib/files.b$;
$trace1$;
intopen_file(stringfile_path,stringmode){intres;
$trace2$;
intm;
$trace3$;
if(mode=="read"){m=0;
$trace4$;
}elseif(mode=="write"){m=1;
$trace5$;
}elseif(mode=="append"){m=2;
$trace6$;
}else{error="no such mode for opening file";
$trace7$;
if($toPanic){print((error+"\n"));
$trace8$;
exit(1);
$trace9$;
};
$trace10$;
returnres;
$trace11$;
};
$trace12$;
res=$open_f(file_path,m);
$trace13$;
returnres;
$trace14$;
};
$trace15$;
stringread_file(intdescriptor_number,intsize){intbytes;
$trace16$;
strings;
$trace17$;
bytes=$read_f(descriptor_number,size,s);
$trace18$;
returns;
$trace19$;
};
$trace20$;
voidclose_file(intdescriptor_number){bool$toPanicOld;
$trace21$;
$toPanicOld=$toPanic;
$trace22$;
$close_f(descriptor_number);
$trace23$;
try{strings;
$trace24$;
intbytes;
$trace25$;
intsize;
$trace26$;
size=1;
$trace27$;
bytes=$read_f(descriptor_number,size,s);
$trace28$;
};
$trace29$;
if(error==""){error=("could not close file with file descriptor number "+str(descriptor_number));
$trace30$;
if($toPanicOld){print((error+"\n"));
$trace31$;
exit(1);
$trace32$;
};
$trace33$;
};
$trace34$;
error="";
$trace35$;
};
$fileprogram.b$;
$trace1$;
voidmain(){intf;
$trace2$;
strings;
$trace3$;
strings2;
$trace4$;
f=open_file("/home/slava/a1.txt","read");
$trace5$;
s=read_file(f,1000);
$trace6$;
s2=read_file(f,1000);
$trace7$;
s=(s+s2);
$trace8$;
print(s);
$trace9$;
print("");
$trace10$;
print("");
$trace11$;
print("");
$trace12$;
close_file(f);
$trace13$;
};
$trace14$;
main();
