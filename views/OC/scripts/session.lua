print("this is my session number 1 ")
eocto.resetWD()
print(string.format("Current Working DIrectory is %s",eocto.getCWD()))
print(string.format("Saved Session %s",eocto.getSession("eocto_cWd")))
eocto.setWD("../../Octopus/views")
print(string.format("Saved Session %s",eocto.getSession("eocto_cWd")))
print(eocto.encodeJSON(eocto.listFiles()))
return 