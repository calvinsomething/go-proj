-- TODO: add preferred language, irl_region with options

CREATE TABLE players (
    ip VARCHAR(51) NOT NULL,
    faction ENUM('H', 'A') NOT NULL,
    race ENUM('dwarf', 'gnome', 'human', 'night elf', 'orc', 'tauren', 'troll', 'undead') NOT NULL,
    class ENUM('druid', 'hunter', 'mage', 'paladin', 'priest', 'rogue', 'shaman', 'warlock', 'warrior') NOT NULL,
    profession1 ENUM('alchemy', 'blacksmithing', 'enchanting', 'engineering', 'herbalism', 'mining', 'tailoring'),
    profession2 ENUM('alchemy', 'blacksmithing', 'enchanting', 'engineering', 'herbalism', 'mining', 'tailoring'),
    weekly_hours INT CHECK (weekly_hours BETWEEN 1 AND 50),
    UNIQUE (ip),
    CHECK (profession2 != profession1)
);

DELIMITER |

CREATE TRIGGER tr_ins_fac BEFORE INSERT ON players
FOR EACH ROW
BEGIN
    IF NEW.race IN ('dwarf', 'gnome', 'human', 'night elf')
    THEN SET NEW.faction = 'A';
    ELSE SET NEW.faction = 'H';
    END IF;
END;|

CREATE TRIGGER tr_up_fac BEFORE UPDATE ON players
FOR EACH ROW
BEGIN
    IF NEW.race IN ('dwarf', 'gnome', 'human', 'night elf')
    THEN SET NEW.faction = 'A';
    ELSE SET NEW.faction = 'H';
    END IF;
END;$$

DELIMITER ;
