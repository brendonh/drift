CubeSet = function(dim) {
    this.dim = dim || 20;
    this.cubes = [];
};

CubeSet.prototype.addAscii = function(lines, colors) {
    var cx, cy;
    for (var i = 0; i < lines.length; i++) {
        var line = lines[i];
        for (var j = 0; j < line.length; j++) {
            if (line[j] == 'X') {
                cx = j;
                cy = i;
                break;
            }
        }
    }

    for (var i = 0; i < lines.length; i++) {
        var line = lines[i];
        for (var j = 0; j < line.length; j++) {
            if (line[j] == ' ') continue;
            this.addCube(j - cx, cy - i, 0, colors[line[j]]);
        }
    }
}

CubeSet.prototype.addCube = function(x, y, z, color) {
    this.cubes.push({'x': x, 'y': y, 'z': z, 'color': color});
}

CubeSet.prototype.toMesh = function() {
    var geom = new THREE.Geometry();
    geom.materials = [];

    var dim = this.dim, dim_half = this.dim / 2;

    var matsByColor = {};

    var byPos = {};
    for (var i = 0; i < this.cubes.length; i++) {
        var cube = this.cubes[i];
        byPos[[cube.x,cube.y,cube.z]] = true;
    }

    for (i = 0; i < this.cubes.length; i++) {
        var cube = this.cubes[i];

        var matIdx = matsByColor[cube.color];
        if (matIdx === undefined) {
            var mat = new THREE.MeshLambertMaterial({color: cube.color, overdraw: true });
            geom.materials.push(mat);
            matIdx = matsByColor = geom.materials.length - 1;
        }

        if (!byPos[[cube.x+1, cube.y, cube.z]])
            this.buildPlane(geom, cube.x, cube.y, cube.z, 'z', 'y', 'x', - 1, - 1, dim,  1, matIdx ); // px

        if (!byPos[[cube.x-1, cube.y, cube.z]])
	        this.buildPlane(geom, cube.x, cube.y, cube.z, 'z', 'y', 'x',   1, - 1, dim, -1, matIdx ); // nx

        if (!byPos[[cube.x, cube.y+1, cube.z]])
	        this.buildPlane(geom, cube.x, cube.y, cube.z, 'x', 'z', 'y',   1,   1, dim,  1, matIdx ); // py

        if (!byPos[[cube.x, cube.y-1, cube.z]])
	        this.buildPlane(geom, cube.x, cube.y, cube.z, 'x', 'z', 'y',   1, - 1, dim, -1, matIdx ); // ny

        if (!byPos[[cube.x, cube.y, cube.z+1]])
	        this.buildPlane(geom, cube.x, cube.y, cube.z, 'x', 'y', 'z',   1, - 1, dim,  1, matIdx ); // pz

        if (!byPos[[cube.x, cube.y, cube.z+-1]])
	       this.buildPlane(geom, cube.x, cube.y, cube.z, 'x', 'y', 'z', - 1, - 1, dim, -1, matIdx ); // nz
    }

	geom.computeCentroids();
	geom.mergeVertices();

    var mesh = new THREE.Mesh( geom, new THREE.MeshFaceMaterial() );

    return mesh;
}

CubeSet.prototype.buildPlane = function(geom, x, y, z, u, v, w, udir, vdir, dim, depth, matIdx ) {
	var dim_half = dim / 2;
    var offset = geom.vertices.length;

	for ( iy = 0; iy < 2; iy ++ ) {
		for ( ix = 0; ix < 2; ix ++ ) {
			var vector = new THREE.Vector3(dim * x, dim * y, dim * z);
			vector[ u ] += ( ix * dim - dim_half ) * udir;
			vector[ v ] += ( iy * dim - dim_half ) * vdir;
			vector[ w ] += dim_half * depth;
			geom.vertices.push( vector );
		}
	}

	var normal = new THREE.Vector3();
	normal[ w ] = depth > 0 ? 1 : - 1;

    var face = new THREE.Face4( 0 + offset, 2 + offset, 3 + offset, 1 + offset );

	face.normal.copy( normal );
    face.vertexNormals.push( normal, normal, normal, normal );
    face.materialIndex = matIdx;

	geom.faces.push( face );
	geom.faceVertexUvs[ 0 ].push( [
		new THREE.UV( 0, 0 ),
		new THREE.UV( 0, 1 ),
		new THREE.UV( 1, 1 ),
		new THREE.UV( 1, 0 )
	] );

}
