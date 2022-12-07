load("../../warpsys.star", "plot")
load("../../warpsys.star", "catalog_input_str")
load("../bootstrap.star", "bootstrap_build_step")
load("../bootstrap.star", "bootstrap_pack_step")

step_build = bootstrap_build_step(
    src=("warpsys.org/gcc", "v11.2.0", "src"),
    extra_inputs = [
        ("warpsys.org/mpfr", "v4.1.0", "src"),
        ("warpsys.org/gmp", "v6.2.1", "src"),
        ("warpsys.org/mpc", "v1.2.1", "src"),
    ],
    script="""cd /src/*
    cp -vpR -v /pkg/warpsys.org/mpfr/* mpfr
    cp -vpR -v /pkg/warpsys.org/gmp/* gmp
    cp -vpR -v /pkg/warpsys.org/mpc/* mpc
    mkdir /prefix/build
    cd /prefix/build
    /src/*/configure --prefix=/warpsys-placeholder-prefix --disable-multilib --enable-languages=c,c++ LDFLAGS=-Wl,-rpath=XORIGIN/../lib
    make
    make DESTDIR=/out install""")

step_pack = bootstrap_pack_step(
    binaries=["gcc"],
    libraries=[
        ("warpsys.org/bootstrap/glibc", "libc.so.6"),
        ("warpsys.org/bootstrap/glibc", "libm.so.6"),
    ])
    
result = plot(steps={"build": step_build, "pack": step_pack})